package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type EnqueueOption func(task *executorTask)

func WithTrackingKey(key string) EnqueueOption {
	return func(task *executorTask) {
		task.key = key
	}
}

type executorTask struct {
	id          uuid.UUID
	key         string
	attempts    int
	maxAttempts int
	payload     WebhookPayload
}

func (e executorTask) shouldTrack() bool {
	return len(e.key) != 0
}

type Executor struct {
	mu     sync.RWMutex
	queues map[string]executorQueue
}

func ProvideExecutor() *Executor {
	return &Executor{
		queues: make(map[string]executorQueue),
	}
}

func (e *Executor) EnqueueEmbed(url string, embed Embed, opts ...EnqueueOption) {
	e.enqueue(url, WebhookPayload{Embeds: []Embed{embed}}, opts...)
}

func (e *Executor) enqueue(url string, payload WebhookPayload, opts ...EnqueueOption) {
	e.mu.Lock()
	queue, ok := e.queues[url]
	if !ok {
		queue = newQueue(url, DefaultTracker())
		e.queues[url] = queue
		go queue.start()
	}
	e.mu.Unlock()

	task := &executorTask{
		id:          uuid.New(),
		attempts:    0,
		maxAttempts: 3,
		payload:     payload,
	}

	for _, opt := range opts {
		opt(task)
	}

	queue.enqueue(task)
}

func newQueue(url string, tracker *Tracker) executorQueue {
	return executorQueue{
		url:     url,
		queue:   make(chan *executorTask),
		tracker: tracker,
	}
}

type executorQueue struct {
	url     string
	queue   chan *executorTask
	tracker *Tracker
}

func (q executorQueue) enqueue(task *executorTask) {
	q.queue <- task
	log.Debug().
		Str("task", task.id.String()).
		Msgf("[Discord] Enqueued task for %s, has %d pending tasks", q.url, len(q.queue))
}

func (q executorQueue) start() {
	log.Info().Msgf("[Discord] Started queue for %s", q.url)

	for {
		select {
		case task := <-q.queue:
			q.processTask(task)
		}
	}
}

func (q executorQueue) processTask(task *executorTask) {
	task.attempts++

	log.Debug().
		Str("task", task.id.String()).
		Int("attempt", task.attempts).
		Msg("[Discord] Started processing task")

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(&task.payload)
	if err != nil {
		log.Error().
			Err(err).
			Str("task", task.id.String()).
			Int("attempt", task.attempts).
			Msg("[Discord] Failed to encode payload")
		return
	}

	var res *http.Response
	if task.shouldTrack() {
		if msgId, tracked := q.tracker.GetMessageID(task.key); tracked {
			req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/messages/%s?wait=true", q.url, msgId), &buf)
			if err != nil {
				log.Error().
					Err(err).
					Str("task", task.id.String()).
					Int("attempt", task.attempts).
					Msg("[Discord] Failed to construct patch request")
				return
			}
			req.Header.Set("Content-Type", "application/json")
			res, err = http.DefaultClient.Do(req)
			if err != nil {
				log.Error().
					Err(err).
					Str("task", task.id.String()).
					Int("attempt", task.attempts).
					Msg("[Discord] Failed to send patch payload")
				return
			}
			goto processRes
		}
	}
	res, err = http.Post(fmt.Sprintf("%s?wait=true", q.url), "application/json", &buf)
	if err != nil {
		log.Error().
			Err(err).
			Str("task", task.id.String()).
			Int("attempt", task.attempts).
			Msg("[Discord] Failed to send payload")
		return
	}

processRes:
	defer res.Body.Close()

	if !isSuccessHttpCode(res.StatusCode) {
		if res.StatusCode == http.StatusTooManyRequests {
			resetTime, err := getRateLimitResetTime(res)
			if err != nil {
				log.Error().Err(err).Msg("[Discord] Got rate limited and expected rate limit headers to be set")
				return
			}
			resetAfter := resetTime.Sub(time.Now())
			log.Warn().
				Str("task", task.id.String()).
				Int("attempt", task.attempts).
				Msgf("[Discord] Got rate limited, re-queueing in %s", resetAfter)

			if resetAfter > 0 {
				time.Sleep(resetAfter)
			}
			q.enqueue(task)
		} else {
			log.Error().
				Str("task", task.id.String()).
				Int("attempt", task.attempts).
				Msgf("[Discord] Received unexpected response from server: %d", res.StatusCode)
			return
		}
	}

	var msg Message
	if err = json.NewDecoder(res.Body).Decode(&msg); err != nil {
		log.Error().
			Err(err).
			Str("task", task.id.String()).
			Int("attempt", task.attempts).
			Msg("[Discord] Failed to parse response body")
		return
	}
	if task.shouldTrack() {
		q.tracker.TrackMessageID(task.key, msg.ID)
	}

	log.Debug().
		Str("task", task.id.String()).
		Int("attempt", task.attempts).
		Msg("[Discord] Successfully processed task")
}

func getRateLimitResetTime(res *http.Response) (time.Time, error) {
	resetTimeStr := res.Header.Get("x-ratelimit-reset")
	if len(resetTimeStr) == 0 {
		return time.Time{}, errors.New("expected an x-ratelimit-reset header")
	}

	resetTime, err := strconv.ParseInt(resetTimeStr, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(resetTime, 0), nil
}

func isSuccessHttpCode(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

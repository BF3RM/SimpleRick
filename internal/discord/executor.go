package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type executorTask struct {
	id          uuid.UUID
	attempts    int
	maxAttempts int
	payload     WebhookPayload
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

func (e *Executor) EnqueueEmbeds(url string, embeds ...Embed) {
	e.enqueue(url, WebhookPayload{Embeds: embeds})
}

func (e *Executor) enqueue(url string, payload WebhookPayload) {
	e.mu.Lock()
	queue, ok := e.queues[url]
	if !ok {
		queue = executorQueue{
			url:   url,
			queue: make(chan *executorTask),
		}
		e.queues[url] = queue
		go queue.start()
	}
	e.mu.Unlock()

	queue.enqueue(&executorTask{
		id:          uuid.New(),
		attempts:    0,
		maxAttempts: 3,
		payload:     payload,
	})
}

type executorQueue struct {
	url   string
	queue chan *executorTask
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

	res, err := http.Post(q.url, "application/json", &buf)
	if err != nil {
		log.Error().
			Err(err).
			Str("task", task.id.String()).
			Int("attempt", task.attempts).
			Msg("[Discord] Failed to send payload")
		return
	}
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

package discord

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"net/http"
	"sync"
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
		Msgf("[Discord] Enqueued new task for %s, has %d pending tasks", q.url, len(q.queue))
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

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		log.Error().
			Err(err).
			Str("task", task.id.String()).
			Int("attempt", task.attempts).
			Msgf("[Discord] Received unexpected response from server: %d", res.StatusCode)
	}

	// TODO: Handle rate limiting

	log.Debug().Str("task", task.id.String()).Msg("Successfully processed task")
}

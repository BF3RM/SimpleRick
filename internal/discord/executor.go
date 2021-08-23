package discord

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"log"
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

func NewExecutor() *Executor {
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
		go queue.process()
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
}

func (q executorQueue) process() {
	log.Printf("Started goroutine to process queue for %s\n", q.url)

	for {
		select {
		case task := <-q.queue:
			q.processTask(task)
		}
	}
}

func (q executorQueue) processTask(task *executorTask) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(&task.payload)
	if err != nil {
		log.Printf("Failed to encode payload of task %s, err=%v\n", task.id, err)
		return
	}

	res, err := http.Post(q.url, "application/json", &buf)
	if err != nil {
		log.Printf("Failed to send payload of task %s, err=%v\n", task.id, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Printf("Unexpected response received for task %s, server responded with %s", task.id, res.StatusCode)
	}
	// TODO: Handle rate limiting
}

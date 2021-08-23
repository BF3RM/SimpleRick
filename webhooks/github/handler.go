package github

import (
	"github.com/google/go-github/github"
	"net/http"
	"simplerick/internal/discord"
)

type WebhookHandler struct {
	executor   *discord.Executor
	secret     []byte
	webhookUri string
}

func New(executor *discord.Executor, secret []byte, webhookUri string) WebhookHandler {
	return WebhookHandler{
		executor:   executor,
		secret:     secret,
		webhookUri: webhookUri,
	}
}

func (h WebhookHandler) Handler(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, h.secret)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	defer r.Body.Close()

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch e := event.(type) {
	case *github.PushEvent:
		err = h.handlePushEvent(e)
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

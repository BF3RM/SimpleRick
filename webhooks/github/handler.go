package github

import (
	"github.com/getsentry/sentry-go"
	"github.com/google/go-github/github"
	"github.com/rs/zerolog/log"
	"net/http"
	"simplerick/internal"
	"simplerick/internal/discord"
)

type WebhookHandler struct {
	executor *discord.Executor
	config   internal.GithubWebhookConfig
}

func ProvideWebhookHandler(executor *discord.Executor, config internal.GithubWebhookConfig) WebhookHandler {
	return WebhookHandler{
		executor: executor,
		config:   config,
	}
}

func (h WebhookHandler) Handler(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, h.config.Secret)
	if err != nil {
		log.Error().Err(err).Msg("[GitHub] Failed to validate payload")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	defer r.Body.Close()

	if e := log.Debug(); e.Enabled() {
		e.Str("body", string(payload)).Msgf("[GitHub] Incoming call")
	}
	sentry.AddBreadcrumb(&sentry.Breadcrumb{
		Level:    sentry.LevelInfo,
		Category: "github",
		Message:  "Incoming webhook call",
		Data: map[string]interface{}{
			"remote": r.RemoteAddr,
		},
	})

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Error().Err(err).Msg("[GitHub] Failed to parse payload")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch e := event.(type) {
	case *github.PushEvent:
		err = h.handlePushEvent(e)
	case *github.CreateEvent:
		err = h.handleCreateEvent(e)
	case *github.DeleteEvent:
		err = h.handleDeleteEvent(e)
	}

	if err != nil {
		log.Error().Err(err).Msg("[GitHub] Failed to process payload")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

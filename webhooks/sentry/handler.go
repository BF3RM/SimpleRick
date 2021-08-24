package sentry

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"simplerick/internal"
	"simplerick/internal/discord"
	"simplerick/internal/sentry"
	"time"
)

type WebhookHandler struct {
	executor *discord.Executor
	config   internal.SentryWebhookConfig
}

func ProvideWebhookHandler(executor *discord.Executor, config internal.SentryWebhookConfig) WebhookHandler {
	return WebhookHandler{
		executor: executor,
		config:   config,
	}
}

func (h WebhookHandler) Handler(w http.ResponseWriter, req *http.Request) {
	payload, err := sentry.ValidatePayload(req, h.config.Secret)
	if err != nil {
		log.Error().Err(err).Msg("[Sentry] Failed to validate payload")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	defer req.Body.Close()

	if e := log.Debug(); e.Enabled() {
		e.Str("body", string(payload)).Msgf("[GitHub] Incoming call")
	}

	action, event, err := sentry.ParseWebhook(sentry.WebhookResource(req), payload)
	if err != nil {
		log.Error().Err(err).Msg("[Sentry] Failed to parse payload")
		w.WriteHeader(http.StatusBadRequest)
	}

	switch e := event.(type) {
	case *sentry.IssueData:
		err = h.handleIssue(action, e)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func (h WebhookHandler) handleIssue(action sentry.EventAction, data *sentry.IssueData) error {

	// TODO: Add support for solving as well
	if action != sentry.IssueCreatedAction && action != sentry.IssueResolvedAction {
		return nil
	}

	builder := discord.NewEmbedBuilder().
		SetTitle(data.Issue.ShortId).
		SetURL(fmt.Sprintf("https://sentry.io/organizations/realitymod-dev-team/issues/%s", data.Issue.Id)).
		SetAuthor(data.Issue.Project.Slug,
			discord.WithAuthorUrl(fmt.Sprintf("https://sentry.io/organizations/realitymod-dev-team/projects/%s", data.Issue.Project.Slug))).
		SetDescription(data.Issue.Title).
		AddField("Status", data.Issue.Status, discord.WithFieldInline()).
		AddField("Level", data.Issue.Level, discord.WithFieldInline()).
		AddField("First Seen", data.Issue.FirstSeen.Format(time.RFC3339)).
		SetFooter("Simple Rick - Sentry").
		AddTimestamp()

	if action == sentry.IssueCreatedAction {
		builder.SetColor(0xE74C3C) // unsolved
	} else {
		builder.SetColor(0x2ECC71) // resolved
	}

	h.executor.EnqueueEmbed(h.config.IssuesWebhookUrl, builder.Build(), discord.WithTrackingKey(data.Issue.Id))

	return nil
}

package github

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/google/go-github/github"
	"github.com/rs/zerolog/log"
	"simplerick/internal/discord"
)

func (h WebhookHandler) handleCreateEvent(event *github.CreateEvent) error {
	if *event.Sender.Type == "Bot" {
		log.Debug().Msg("[GitHub] Ignored create event from bot")
		return nil
	}

	if *event.RefType != "branch" {
		return nil
	}

	sentry.AddBreadcrumb(&sentry.Breadcrumb{
		Category: "github",
		Message:  "Handling create event",
		Data: map[string]interface{}{
			"repo":   *event.Repo.Name,
			"sender": *event.Sender.Login,
			"ref":    *event.Ref,
		},
		Level: sentry.LevelInfo,
	})

	builder := discord.NewEmbedBuilder().
		SetColor(0x00BCD4).
		SetAuthor(*event.Sender.Login,
			discord.WithAuthorUrl(*event.Sender.HTMLURL),
			discord.WithAuthorIcon(*event.Sender.AvatarURL)).
		SetURL(fmt.Sprintf("%s/tree/%s", *event.Repo.HTMLURL, *event.Ref)).
		SetDescription(fmt.Sprintf("Created branch **%s** on **%s**", *event.Ref, *event.Repo.Name)).
		SetFooter("Simple Rick - GitHub").
		AddTimestamp()

	h.executor.EnqueueEmbed(h.config.ChangelogWebhookUrl, builder.Build())

	return nil
}

func (h WebhookHandler) handleDeleteEvent(event *github.DeleteEvent) error {
	if *event.Sender.Type == "Bot" {
		log.Debug().Msg("[GitHub] Ignored delete event from bot")
		return nil
	}

	if *event.RefType != "branch" {
		return nil
	}

	sentry.AddBreadcrumb(&sentry.Breadcrumb{
		Category: "github",
		Message:  "Handling delete event",
		Data: map[string]interface{}{
			"repo":   *event.Repo.Name,
			"sender": *event.Sender.Login,
			"ref":    *event.Ref,
		},
		Level: sentry.LevelInfo,
	})

	builder := discord.NewEmbedBuilder().
		SetColor(0x00BCD4).
		SetAuthor(*event.Sender.Login,
			discord.WithAuthorUrl(*event.Sender.HTMLURL),
			discord.WithAuthorIcon(*event.Sender.AvatarURL)).
		SetDescription(fmt.Sprintf("Deleted branch **%s** of **%s**", *event.Ref, *event.Repo.Name)).
		SetFooter("Simple Rick - GitHub").
		AddTimestamp()

	h.executor.EnqueueEmbed(h.config.ChangelogWebhookUrl, builder.Build())

	return nil
}

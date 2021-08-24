package github

import (
	"fmt"
	"github.com/google/go-github/github"
	"simplerick/internal/discord"
)

func (h WebhookHandler) handleCreateEvent(event *github.CreateEvent) error {
	if *event.RefType != "branch" {
		return nil
	}

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
	if *event.RefType != "branch" {
		return nil
	}

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

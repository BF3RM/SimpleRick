package github

import (
	"fmt"
	"github.com/google/go-github/github"
	"simplerick/internal/discord"
	"simplerick/internal/utils"
	"strings"
)

func (h WebhookHandler) handlePushEvent(event *github.PushEvent) error {
	branch := (*event.Ref)[len("refs/heads/"):]
	lenCommits := len(event.Commits)

	if lenCommits == 0 {
		return nil
	}

	builder := discord.NewEmbedBuilder().
		SetColor(0x00BCD4).
		SetAuthor(*event.Sender.Login,
			discord.WithAuthorUrl(*event.Sender.HTMLURL),
			discord.WithAuthorIcon(*event.Sender.AvatarURL)).
		SetDescription(fmt.Sprintf("to branch **%s** of **%s**", branch, *event.Repo.Name)).
		SetFooter("Simple Rick - GitHub").
		AddTimestamp()

	if lenCommits == 1 {
		builder.
			SetTitle("Pushed a commit").
			SetURL(*event.Commits[0].URL)
	} else {
		builder.
			SetTitle(fmt.Sprintf("Pushed %d commits", lenCommits)).
			SetURL(*event.Compare)
	}

	// If more than 25 commits, grab last 25
	if lenCommits > 25 {
		event.Commits = event.Commits[lenCommits-26 : lenCommits-1]
	}

	for _, commit := range event.Commits {
		sha := (*commit.ID)[:7]
		messages := strings.Split(*commit.Message, "\n")

		title := fmt.Sprintf("`%s` %s", sha, messages[0])
		description := fmt.Sprintf("- **%s**", *commit.Committer.Name)

		if len(messages) > 1 {
			description = utils.Ellipsis(strings.Join(messages[1:], "\n"), 255-len(description)) + "\n" + description
		}

		builder.AddField(title, description)
	}

	h.executor.EnqueueEmbeds(h.config.ChangelogWebhookUrl, builder.Build())

	return nil
}

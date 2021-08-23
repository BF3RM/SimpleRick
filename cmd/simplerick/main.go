package main

import (
	"net/http"
	"simplerick/internal/discord"
	"simplerick/internal/env"
	"simplerick/webhooks/github"
	"simplerick/webhooks/sentry"
)

func main() {
	executor := discord.NewExecutor()

	githubWebhook := github.New(executor,
		[]byte(env.GetString("GITHUB_WEBHOOK_SECRET", "")),
		env.GetString("DISCORD_CHANGELOG_WEBHOOK_URL", ""))

	sentryWebhook := sentry.New(executor,
		[]byte(env.GetString("SENTRY_WEBHOOK_SECRET", "")),
		env.GetString("DISCORD_ISSUES_WEBHOOK_URL", ""))

	http.HandleFunc("/api/v1/webhooks/github", githubWebhook.Handler)
	http.HandleFunc("/api/v1/webhooks/sentry", sentryWebhook.Handler)
	http.ListenAndServe(":3000", nil)
}

package internal

import (
	"errors"
	"simplerick/internal/env"
)

type SentryWebhookConfig struct {
	Secret           []byte
	IssuesWebhookUrl string
}

type GithubWebhookConfig struct {
	Secret              []byte
	ChangelogWebhookUrl string
	ReleasesWebhookUrl  string
}

func ProvideSentryWebhookConfig() (SentryWebhookConfig, error) {
	secret := env.GetBytes("SENTRY_WEBHOOK_SECRET", nil)
	issuesWebhookUrl := env.GetString("SENTRY_ISSUES_WEBHOOK_URL", "")

	if len(issuesWebhookUrl) == 0 {
		return SentryWebhookConfig{}, errors.New("environment variable SENTRY_ISSUES_WEBHOOK_URL is not set")
	}

	return SentryWebhookConfig{secret, issuesWebhookUrl}, nil
}

func ProvideGithubWebhookConfig() (GithubWebhookConfig, error) {
	secret := env.GetBytes("GITHUB_WEBHOOK_SECRET", nil)
	changelogWebhookUrl := env.GetString("GITHUB_CHANGES_WEBHOOK_URL", "")
	releasesWebhookUrl := env.GetString("GITHUB_RELEASES_WEBHOOK_URL", "")

	if len(changelogWebhookUrl) == 0 {
		return GithubWebhookConfig{}, errors.New("environment variable GITHUB_CHANGES_WEBHOOK_URL is not set")
	}

	if len(releasesWebhookUrl) == 0 {
		return GithubWebhookConfig{}, errors.New("environment variable GITHUB_RELEASES_WEBHOOK_URL is not set")
	}

	return GithubWebhookConfig{secret, changelogWebhookUrl, releasesWebhookUrl}, nil
}

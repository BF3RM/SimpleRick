package internal

import (
	"github.com/google/wire"
	"simplerick/internal/discord"
)

var Set = wire.NewSet(ProvideGithubWebhookConfig, ProvideSentryWebhookConfig, discord.ProvideExecutor)

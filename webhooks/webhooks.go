package webhooks

import (
	"github.com/google/wire"
	"simplerick/webhooks/github"
	"simplerick/webhooks/sentry"
)

var Set = wire.NewSet(github.ProvideWebhookHandler, sentry.ProvideWebhookHandler)

package main

import (
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/google/wire"
	"github.com/gorilla/mux"
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"net/http"
	"path"
	"simplerick/internal/env"
	"simplerick/internal/logging"
	github_webhook "simplerick/webhooks/github"
	sentry_webhook "simplerick/webhooks/sentry"
	"time"
)

func setupLogger() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	writers := []io.Writer{
		zerolog.ConsoleWriter{Out: colorable.NewColorableStdout()},
		&logging.SentryWriter{},
		&lumberjack.Logger{
			Filename:   path.Join("logs", "simplerick.log"),
			MaxSize:    100,
			MaxAge:     7,
			MaxBackups: 3,
		},
	}
	log.Logger = zerolog.New(zerolog.MultiLevelWriter(writers...)).With().Timestamp().Logger()
}

func setupSentry() bool {
	sentryDsn := env.GetString("SENTRY_DSN", "")
	if sentryDsn == "" {
		log.Debug().Msg("[Main] Sentry disabled")
		return false
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:         sentryDsn,
		Release:     "simplerick@0.1.0",
		Environment: "development",
	})

	if err != nil {
		log.Fatal().Err(err).Msg("Sentry init failed")
	}

	return true
}

func main() {
	setupLogger()

	if setupSentry() {
		defer sentry.Flush(2 * time.Second)
	}

	ctx := context.Background()
	app, err := setupApplication(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("[Main] Failed to set up application")
	}

	if err = app.Start(); err != nil {
		log.Fatal().Err(err).Msg("[Main] Failed to start application")
	}
}

var applicationSet = wire.NewSet(
	newApplication,
	newRouter,
	wire.Bind(new(http.Handler), new(*mux.Router)),
)

func newRouter(githubWebhook github_webhook.WebhookHandler, sentryWebhook sentry_webhook.WebhookHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/webhooks/github", githubWebhook.Handler)
	r.HandleFunc("/api/v1/webhooks/sentry", sentryWebhook.Handler)
	return r
}

func newApplication(handler http.Handler) application {
	return application{handler}
}

type application struct {
	handler http.Handler
}

func (app application) Start() error {
	log.Info().Msg("[Main] Listening to port 3000")
	return http.ListenAndServe(":3000", app.handler)
}

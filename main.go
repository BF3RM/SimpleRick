package main

import (
	"context"
	"github.com/google/wire"
	"github.com/gorilla/mux"
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"net/http"
	"path"
	"simplerick/webhooks/github"
	"simplerick/webhooks/sentry"
)

func setupLogger() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	var writers []io.Writer

	writers = append(writers, zerolog.ConsoleWriter{Out: colorable.NewColorableStdout()})
	writers = append(writers, &lumberjack.Logger{
		Filename:   path.Join("logs", "simplerick.log"),
		MaxSize:    100,
		MaxAge:     7,
		MaxBackups: 3,
	})

	log.Logger = zerolog.New(zerolog.MultiLevelWriter(writers...)).With().Timestamp().Logger()
}

func main() {
	setupLogger()

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

func newRouter(githubWebhook github.WebhookHandler, sentryWebhook sentry.WebhookHandler) *mux.Router {
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

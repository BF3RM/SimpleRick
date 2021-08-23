package main

import (
	"context"
	"github.com/google/wire"
	"github.com/gorilla/mux"
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"simplerick/webhooks/github"
	"simplerick/webhooks/sentry"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: colorable.NewColorableStdout()})

	zerolog.SetGlobalLevel(zerolog.DebugLevel)

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

// +build wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"simplerick/internal"
	"simplerick/webhooks"
)

func setupApplication(ctx context.Context) (application, error) {
	wire.Build(
		internal.Set,
		webhooks.Set,
		applicationSet,
	)
	return application{}, nil
}

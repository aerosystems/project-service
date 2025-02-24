package main

import (
	"context"
	"golang.org/x/sync/errgroup"
)

// @title Project Service
// @version 2.0.0
// @description A part of microservice infrastructure, who responsible for managing user Projects

// @contact.name Artem Kostenko
// @contact.url https://github.com/aerosystems

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey X-Api-Key
// @in header
// @name X-Api-Key
// @description Should contain Token, digits and letters, 64 symbols length

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Should contain Access JWT Token, with the Bearer started

// @host gw.verifire.dev/project
// @schemes https
// @BasePath /
func main() {
	app := InitApp()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)

	switch app.cfg.Proto {
	case "http":
		group.Go(func() error {
			return app.httpServer.Run()
		})
	case "grpc":
		group.Go(func() error {
			return app.grpcServer.Run()
		})
	default:
		app.log.Fatalf("unknown protocol: %s", app.cfg.Proto)
	}

	group.Go(func() error {
		return app.handleSignals(ctx, cancel)
	})

	if err := group.Wait(); err != nil {
		app.log.Errorf("error occurred: %v", err)
	}
}

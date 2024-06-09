package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func (app *App) handleSignals(ctx context.Context, cancel context.CancelFunc) error {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	select {
	case <-signalCh:
		// Handle graceful shutdown
		return app.gracefulShutdown(ctx, cancel)
	case <-ctx.Done():
		// Context cancelled, shutdown initiated elsewhere
		return nil
	}
}

func (app *App) gracefulShutdown(_ context.Context, cancel context.CancelFunc) error {
	cancel()
	app.log.Fatalf("app is shutting down")
	return nil
}

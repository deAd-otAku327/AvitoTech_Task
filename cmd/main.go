package main

import (
	"errors"
	"log"
	"log/slog"
	"merch_shop/cmd/app"
	"merch_shop/internal/config"
	"merch_shop/internal/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const configPath = "configs/local.yaml"

func main() {
	cfg, err := config.New(configPath)
	reportOnError(err)

	app, err := app.New(cfg, logger.New(os.Stdout))
	reportOnError(err)

	go func() {
		err = app.Run()
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Printf("server error: %v", err)
			}
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch
	slog.Info("starting server shutdown")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
	slog.Info("server shutdown completed")
}

func reportOnError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

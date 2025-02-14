package main

import (
	"errors"
	"log"
	"merch_shop/cmd/app"
	"merch_shop/internal/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const configPath = "../configs/dev.yaml"

func main() {
	cfg, err := config.New(configPath)
	reportOnError(err)

	app, err := app.New(cfg)
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

	select {
	case <-ch:
		log.Println("starting server shutdown")
		if err := app.Shutdown(); err != nil {
			log.Fatalf("server shutdown failed: %v", err)
		}
		log.Println("server shutdown completed")
	}
}

func reportOnError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

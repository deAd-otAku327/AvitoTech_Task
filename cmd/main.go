package main

import (
	"log"
	"merch_shop/cmd/app"
	"merch_shop/internal/config"
)

const configPath = "../configs/dev.yaml"

func main() {
	cfg, err := config.New(configPath)
	reportOnError(err)

	app, err := app.New(cfg)
	reportOnError(err)

	err = app.Run()
	reportOnError(err)
}

func reportOnError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

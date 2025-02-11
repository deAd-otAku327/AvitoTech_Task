package app

import (
	"merch_shop/internal/config"
	"net/http"
)

type App struct {
	cfg    *config.Config
	server *http.Server
}

func New()

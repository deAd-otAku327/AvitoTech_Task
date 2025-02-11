package app

import (
	"fmt"
	"log"
	"merch_shop/internal/config"
	"merch_shop/internal/db"
	"merch_shop/pkg/tokenizer"
	"net/http"

	"github.com/gorilla/mux"
)

const AppName string = "merch-shop service"

type App struct {
	cfg    *config.Config
	server *http.Server
}

func New(cfg *config.Config) (*App, error) {
	tokenizer := tokenizer.New(AppName, cfg.SecretKey)
	storage, err := db.New(cfg.DB)
	if err != nil {
		return nil, err
	}

	router := mux.NewRouter()

	return &App{
		cfg: cfg,
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
			Handler: router,
		},
	}, nil
}

func (app *App) Run() error {
	log.Printf("server starting on %s", app.server.Addr)
	return app.server.ListenAndServe()
}

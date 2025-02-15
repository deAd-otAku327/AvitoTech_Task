package app

import (
	"context"
	"fmt"
	"log/slog"
	"merch_shop/internal/config"
	"merch_shop/internal/db"
	"merch_shop/internal/handlers"
	"merch_shop/internal/service"
	"merch_shop/pkg/cryptor"
	"merch_shop/pkg/middleware"
	"merch_shop/pkg/tokenizer"
	"net/http"

	"github.com/gorilla/mux"
)

const AppName string = "merch-shop service"

type App struct {
	cfg    *config.Config
	server *http.Server
}

func New(cfg *config.Config, logger *slog.Logger) (*App, error) {
	tokenizer := tokenizer.New(AppName, cfg.SecretKey)
	cryptor := cryptor.New()

	storage, err := db.New(cfg.DB)
	if err != nil {
		return nil, err
	}

	service := service.New(storage, logger, cryptor, tokenizer)

	controller := handlers.New(service)

	router := mux.NewRouter()
	router.Use(middleware.RpsLimit(cfg.RPS))
	router.Use(middleware.ResponseTimeLimit(cfg.ResponseTime))
	router.Use(middleware.Logging(logger))

	authRouter := router.PathPrefix("/api/auth").Subrouter()
	authRouter.HandleFunc("", controller.Auth()).Methods(http.MethodPost)

	businessRouter := router.PathPrefix("/api").Subrouter()
	businessRouter.Use(middleware.Auth(tokenizer))

	businessRouter.HandleFunc("/info", controller.GetInfo()).Methods(http.MethodGet)
	businessRouter.HandleFunc("/buy/{item:[0-9]+}", controller.BuyItem()).Methods(http.MethodGet)
	businessRouter.HandleFunc("/sendCoin", controller.SendCoin()).Methods(http.MethodPost)

	return &App{
		cfg: cfg,
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
			Handler: router,
		},
	}, nil
}

func (app *App) Run() error {
	slog.Info("server starting on", slog.String("address", app.server.Addr))
	return app.server.ListenAndServe()
}

func (app *App) Shutdown() error {
	return app.server.Shutdown(context.Background())
}

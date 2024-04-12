package app

import "github.com/mzhn-sochi/auth-service/internal/config"

type App struct {
	cfg *config.Config
}

func newApp(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

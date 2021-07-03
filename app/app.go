package app

import (
	"database/sql"
	"fmt"

	"github.com/dhable/mini-auth/jwt"
	"github.com/gorilla/mux"
)

type App struct {
	Router   *mux.Router
	Db       *sql.DB
	KeyStore *jwt.KeyStore
}

func NewApp() (*App, error) {
	app := &App{}

	app.initializeRoutes()

	if err := app.initializeDatabase(); err != nil {
		return nil, fmt.Errorf("fatal DB init error: %s", err)
	}

	if err := app.initializeKeyStore(); err != nil {
		return nil, fmt.Errorf("fatal KeyStore init error: %s", err)
	}

	return app, nil
}

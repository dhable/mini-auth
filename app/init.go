package app

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dhable/mini-auth/jwt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/gorilla/mux"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func (app *App) initializeRoutes() {
	if app.Router == nil {
		app.Router = mux.NewRouter()
		app.Router.HandleFunc("/authenticate", app.authenticateUser).Methods("POST")
		app.Router.HandleFunc("/profile", app.currentUserProfile).Methods("GET")
		app.Router.HandleFunc("/jwks", app.jwtPublicKeySet).Methods("GET")
	} else {
		fmt.Println("skipping init router")
	}
}

func (app *App) initializeDatabase() error {
	if app.Db == nil {
		db, err := sql.Open("sqlite3", "./mini-auth.db")
		if err != nil {
			return err
		}
		app.Db = db
	}

	driver, err := sqlite3.WithInstance(app.Db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "sqlite3", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err.Error() == "no change" {
		return nil
	} else {
		return err
	}
}

func (app *App) initializeKeyStore() error {
	ks, err := jwt.SingleKey("mini-auth", time.Hour)
	if err != nil {
		return err
	}

	app.KeyStore = ks
	return nil
}

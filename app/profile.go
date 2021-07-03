package app

import (
	"net/http"
	"strings"

	"github.com/dhable/mini-auth/models"
)

func (app *App) currentUserProfile(rw http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		jwt := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := app.KeyStore.ParseJWT(jwt)
		if err == nil {
			rows, err := app.Db.Query("SELECT email, name, location FROM users WHERE email = ?", claims.Email)
			if rows != nil {
				defer rows.Close()
			}

			if err == nil && rows.Next() {
				var (
					email    string
					name     string
					location string
				)

				if err := rows.Scan(&email, &name, &location); err == nil {
					ok(rw, models.UserProfile{
						Email:    email,
						Name:     name,
						Location: location,
					})
					return
				}
			}
		}
	}

	unauthorized(rw)
}

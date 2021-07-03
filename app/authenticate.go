package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dhable/mini-auth/models"
	"golang.org/x/crypto/bcrypt"
)

func parseAuthRequest(r *http.Request) (*models.AuthRequest, error) {
	defer r.Body.Close()

	var authReq models.AuthRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&authReq); err != nil {
		log.Fatalf("failed to parse request body: %s", err)
		return nil, err
	} else {
		return &authReq, nil
	}
}

func (app *App) authenticateUser(rw http.ResponseWriter, r *http.Request) {
	authReq, err := parseAuthRequest(r)
	if err != nil {
		badRequest(rw, err)
		return
	}

	rows, err := app.Db.Query("SELECT email, hashedpw FROM users WHERE email = ?", authReq.Email)
	if rows != nil {
		defer rows.Close()
	}

	if err == nil && rows.Next() {
		var (
			email          string
			hashedPassword []byte
		)

		if err := rows.Scan(&email, &hashedPassword); err == nil {
			if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(authReq.Password)); err == nil {
				jwt, err := app.KeyStore.GenerateJWT(email)
				if err == nil {
					ok(rw, models.AuthResponse{
						Jwt: jwt,
					})
					return
				}
			}
		}
	}

	unauthorized(rw)
}

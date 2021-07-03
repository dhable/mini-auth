package app

import (
	"net/http"
)

func (app *App) jwtPublicKeySet(rw http.ResponseWriter, r *http.Request) {
	jwks := app.KeyStore.JWTS()
	ok(rw, jwks)
}

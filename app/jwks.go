package app

import (
	"net/http"
)

func (app *App) JwtPublicKeySet(rw http.ResponseWriter, r *http.Request) {
	jwks := app.KeyStore.JWTS()
	ok(rw, jwks)
}

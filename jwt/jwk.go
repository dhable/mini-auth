package jwt

import (
	"github.com/square/go-jose/v3"
)

func (ks *KeyStore) JWTS() jose.JSONWebKeySet {

	jwkArr := make([]jose.JSONWebKey, 0)
	for _, priKey := range ks.allKeys() {
		jwkArr = append(jwkArr, jose.JSONWebKey{
			Key:       priKey,
			Algorithm: "RS256",
			Use:       "sig",
		})
	}

	return jose.JSONWebKeySet{
		Keys: jwkArr,
	}
}

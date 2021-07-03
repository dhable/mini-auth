package jwt

import (
	"fmt"
	"time"

	"github.com/square/go-jose/v3"
	"github.com/square/go-jose/v3/jwt"
)

type UserClaims struct {
	jwt.Claims
	Email string `json:"email"`
}

func (ks *KeyStore) GenerateJWT(email string) (string, error) {
	expireAt := time.Now().Add(ks.tokenTTL)
	claims := UserClaims{
		Email: email,
		Claims: jwt.Claims{
			Issuer: ks.issuer,
			Expiry: jwt.NewNumericDate(expireAt),
		},
	}

	opts := jose.SignerOptions{}
	opts.WithType("JWT")

	signKey := jose.SigningKey{
		Algorithm: jose.RS256,
		Key:       ks.currentKey(),
	}

	signer, err := jose.NewSigner(signKey, &opts)
	if err != nil {
		return "", err
	}

	return jwt.Signed(signer).Claims(claims).CompactSerialize()
}

func (ks *KeyStore) ParseJWT(signedJwt string) (*UserClaims, error) {
	token, err := jwt.ParseSigned(signedJwt)
	if err != nil {
		return nil, err
	}

	var claims *UserClaims
	for _, priKey := range ks.allKeys() {
		pubKey := &priKey.PublicKey
		claim := new(UserClaims)
		if err := token.Claims(pubKey, claim); err == nil {
			claims = claim
			break
		}
	}

	if claims == nil {
		return nil, fmt.Errorf("failed to find signing key")
	}

	expected := jwt.Expected{
		Issuer: ks.issuer,
		Time:   time.Now(),
	}

	if err := claims.Validate(expected); err != nil {
		return nil, err
	}

	return claims, nil
}

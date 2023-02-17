package token

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type JWTClaim struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateJWT(username string, secret string, expiry time.Duration) (string, error) {
	expirationTime := time.Now().Add(expiry)
	claims := &JWTClaim{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(signedToken string, secret string) error {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)
	if err != nil {
		return err
	}
	_, ok := token.Claims.(*JWTClaim)
	if !ok {
		return errors.New("couldn't parse claims")
	}

	return nil
}

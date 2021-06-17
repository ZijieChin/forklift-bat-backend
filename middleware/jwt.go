package middleware

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	UserID string
	jwt.StandardClaims
}

var Cypher = []byte("cocacola") // Cypher used by sha-256

const TokenExpireDuration = time.Hour * 8 // Token expires in 8 hours

func GenToken(userid string) (string, error) {
	// Custom claims include user id and expire time
	claims := &CustomClaims{
		UserID: userid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(Cypher)
}

func ParseToken(tokenstring string) (*CustomClaims, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenstring, claims, func(t *jwt.Token) (interface{}, error) {
		return Cypher, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid { // check token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

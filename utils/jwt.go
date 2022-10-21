package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("I'mGoingToBeAGolangDeveloper")

type Claims struct {
	UserID string
	Role   string
	jwt.StandardClaims
}

func Encode(userID string, role string) (tokenString string, err error) {
	tokenExpirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(secretKey)
	if err != nil {
		err = fmt.Errorf("error generating token, err: %v", err)
		return
	}
	return
}

func Decode(token string) (claims *Claims, err error) {
	if token == "" {
		err := fmt.Errorf("authorization token is missing")
		return nil, err
	}

	claims = &Claims{}

	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		err = fmt.Errorf("unauthorized, err: %v", err)
		return
	}
	return
}

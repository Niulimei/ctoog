package utils

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	UserName string
	jwt.StandardClaims
}

func CreateJWT(username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		UserName: username,
	})
	tokenString, err := token.SignedString([]byte("Valar Morghulis"))
	if err != nil {
		return ""
	}
	return tokenString
}

func Verify(tokenString string) (string, bool) {
	if tokenString == "" {
		return "", false
	}
	token, _ := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("Valar Morghulis"), nil
	})
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.UserName, true
	} else {
		return "", false
	}
}

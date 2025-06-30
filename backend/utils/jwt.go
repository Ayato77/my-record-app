package utils

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
    UserID uint
    jwt.RegisteredClaims
}

func GenerateToken(userID uint) (string, error) {
    claims := &Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //NewWithClaims creates a new Token with the specified signing method and claims
    return token.SignedString(jwtKey) //SignedString returns the complete, signed token as a string
}

func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
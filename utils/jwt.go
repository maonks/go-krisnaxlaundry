package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(
	userID int64,
	roleID int64,
) (string, error) {

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userID,
			"role_id": roleID,

			"exp": time.Now().
				Add(24 * time.Hour).
				Unix(),
		},
	)

	return token.SignedString(
		[]byte(os.Getenv("JWT_SECRET")),
	)
}

func ParseJWT(tokenString string) (*jwt.Token, error) {

	return jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {

			return []byte(
				os.Getenv("JWT_SECRET"),
			), nil
		},
	)
}

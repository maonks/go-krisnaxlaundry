package auth

import (
	"context"
	"os"

	"github.com/maonkscode/go-kresnaxlaundry/users"

	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

func LoginWithGoogle(
	db *gorm.DB,
	idToken string,
) (*users.User, error) {

	payload, err := idtoken.Validate(
		context.Background(),
		idToken,
		os.Getenv("GOOGLE_CLIENT_ID"),
	)

	if err != nil {
		return nil, err
	}

	email := payload.Claims["email"].(string)

	var user users.User

	err = db.
		Where("email = ?", email).
		Where("is_active = ?", true).
		First(&user).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

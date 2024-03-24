package token

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/buguzei/medods-task/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

const (
	refreshTokenSize = 32
)

func NewAccessToken(user models.User, secretKey string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)

	claims := token.Claims.(jwt.MapClaims)

	claims["sub"] = user.GUID

	strToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("error of signing: %w", err)
	}

	return strToken, nil
}

func NewRefreshToken() (string, error) {
	rb := make([]byte, refreshTokenSize)

	_, err := rand.Read(rb)
	if err != nil {
		return "", fmt.Errorf("error of reading: %w", err)
	}

	return base64.URLEncoding.EncodeToString(rb), nil
}

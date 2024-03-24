package repo

import (
	"context"
	"github.com/buguzei/medods-task/internal/models"
)

type AuthRepo interface {
	NewRefreshToken(ctx context.Context, user models.User, refreshToken string) error
	UpdateRefreshToken(ctx context.Context, user models.User, refreshToken string) error
	IsUserExist(ctx context.Context, user models.User) (bool, error)
	IsRefreshTokenExists(ctx context.Context, refreshToken string) (bool, string, error)
}

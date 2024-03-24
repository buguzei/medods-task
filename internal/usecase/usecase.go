package usecase

import (
	"context"
	"github.com/buguzei/medods-task/internal/models"
	"github.com/buguzei/medods-task/internal/repo"
	"github.com/buguzei/medods-task/pkg/config"
)

type UseCase struct {
	repo repo.AuthRepo
	cfg  *config.Config
}

func NewUseCase(repo repo.AuthRepo, cfg *config.Config) *UseCase {
	return &UseCase{repo: repo, cfg: cfg}
}

type AuthUC interface {
	NewPair(ctx context.Context, user models.User) (*models.TokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (*models.TokenPair, error)
}

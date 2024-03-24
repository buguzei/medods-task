package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/buguzei/medods-task/internal/models"
	"github.com/buguzei/medods-task/internal/token"
	"golang.org/x/crypto/bcrypt"
)

func (uc UseCase) NewPair(ctx context.Context, user models.User) (*models.TokenPair, error) {
	isUserExist, err := uc.repo.IsUserExist(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error of checking for user existing: %w", err)
	}

	if isUserExist {
		return nil, errors.New("user already exists")
	}

	accessToken, err := token.NewAccessToken(user, uc.cfg.SecretKey)

	refreshToken, err := token.NewRefreshToken()

	hash, err := bcryptHash(refreshToken)
	if err != nil {
		return nil, err
	}

	err = uc.repo.NewRefreshToken(ctx, user, hash)
	if err != nil {
		return nil, fmt.Errorf("error of mongo new refresh: %w", err)
	}

	pair := &models.TokenPair{
		Access:  accessToken,
		Refresh: refreshToken,
	}

	return pair, nil
}

func (uc UseCase) Refresh(ctx context.Context, refreshToken string) (*models.TokenPair, error) {
	var res models.TokenPair

	isRefreshTokenExist, userID, err := uc.repo.IsRefreshTokenExists(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("error of checking for refresh existing: %w", err)
	}

	if !isRefreshTokenExist {
		return nil, errors.New("invalid refresh token")
	}

	user := models.User{GUID: userID}

	accessToken, err := token.NewAccessToken(user, uc.cfg.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("error of getting new access token: %w", err)
	}

	newRefreshToken, err := token.NewRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("error of getting new refresh token: %w", err)
	}

	hashedRefreshToken, err := bcryptHash(newRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("error of bcrypt refresh token: %w", err)
	}

	err = uc.repo.UpdateRefreshToken(ctx, user, hashedRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("error of updating refresh in db: %w", err)
	}

	res.Refresh = newRefreshToken
	res.Access = accessToken

	return &res, nil
}

func bcryptHash(src string) (string, error) {
	bHash, err := bcrypt.GenerateFromPassword([]byte(src), 14)
	if err != nil {
		return "", fmt.Errorf("error of bcrypt: %w", err)
	}

	return string(bHash), nil
}

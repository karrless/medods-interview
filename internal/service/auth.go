// Package service provides tools for working with users
package service

import (
	"context"
	"medods-jwt/internal/models"
	"medods-jwt/pkg/errors"
	"medods-jwt/pkg/jwt"
	"medods-jwt/pkg/utils"
)

// UserRepo is a repository for users
type AuthRepo interface {
	GetUser(guid, refresh *string) (*models.User, error)
	UpdateRefreshToken(guid, refreshToken, publicKey string) error
	CreateRefreshToken(guid, refreshToken, publicKey string) error
}

// EmailRepo is a repository for emails
type EmailRepo interface {
	SendWarning(email, userIP, currentIP string) error
}

// AuthService is a service for users
type AuthService struct {
	ctx       *context.Context
	secretKey string
	repo      AuthRepo
	emailRepo EmailRepo
}

// AuthService returns a new user service
func NewAuthService(ctx *context.Context, secretKey string, repo AuthRepo, emailRepo EmailRepo) *AuthService {
	return &AuthService{ctx: ctx, secretKey: secretKey, repo: repo, emailRepo: emailRepo}
}

// getTokens returns access, refresh tokens and public key
func getTokens(guid, ip string) (string, string, string, error) {
	accessToken := jwt.NewAccessToken(guid, ip, "secret")
	refreshToken, err := jwt.NewRefreshToken()
	if err != nil {
		return "", "", "", errors.ErrCreateRefreshToken
	}
	publicKey := utils.GenerateGUID()

	return accessToken, refreshToken, publicKey, nil
}

// RetryCreateTokens retries creating unique access and refresh tokens
func retryCreateTokens(f func() (string, string, string, error)) (string, string, string, error) {
	var err error
	for true {
		accessToken, refreshToken, publicKey, err := f()
		if err == nil {
			return accessToken, refreshToken, publicKey, nil
		}
		if err != errors.ErrRefreshTokenNotUnique {
			break
		}
	}
	return "", "", "", err
}

// GetAccessToken returns a pair of access and refresh tokens
func (uc *AuthService) GetAccessToken(guid, ip string) (accessToken, refreshToken, publicKey string, err error) {
	if !utils.IsGUID(guid) {
		return "", "", "", errors.ErrInvalidGUID
	}

	user, err := uc.repo.GetUser(&guid, nil)
	if err != nil {
		return "", "", "", err
	}

	if user.IP != ip {
		err = uc.emailRepo.SendWarning(user.Email, user.IP, ip)
		if err != nil {
			return "", "", "", err
		}
	}

	if user.Refresh.Valid {
		return "", "", "", errors.ErrUserAlreadyHasAccessToken
	}

	accessToken, refreshToken, publicKey, err = retryCreateTokens(func() (string, string, string, error) {
		accessToken, refreshToken, publicKey, err := getTokens(guid, ip)
		if err != nil {
			return "", "", "", err
		}
		if err := uc.repo.CreateRefreshToken(user.GUID, refreshToken, publicKey); err != nil {
			return "", "", "", err
		}
		return accessToken, refreshToken, publicKey, err
	})

	encryptedRefreshToken := utils.EncodeBase64(refreshToken)
	return accessToken, encryptedRefreshToken, publicKey, nil
}

// RefreshToken returns a pair of access and refresh tokens
func (uc *AuthService) RefreshToken(refresh, public, ip string) (accessToken, refreshToken, publicKey string, err error) {
	refreshToken, err = utils.DecodeBase64(refresh)
	if err != nil {
		return "", "", "", errors.ErrInvalidRefreshToken
	}

	user, err := uc.repo.GetUser(nil, &public)
	if err != nil {
		return "", "", "", err
	}

	if user.IP != ip {
		err = uc.emailRepo.SendWarning(user.Email, user.IP, ip)
		if err != nil {
			return "", "", "", err
		}
	}
	if user.Refresh.Valid {
		ok, _ := utils.CheckStirngHash(refreshToken, user.Refresh.String)
		if !ok {
			return "", "", "", errors.ErrInvalidRefreshToken
		}
	}

	accessToken, refreshToken, publicKey, err = retryCreateTokens(func() (string, string, string, error) {
		accessToken, refreshToken, publicKey, err := getTokens(user.GUID, ip)
		if err != nil {
			return "", "", "", err
		}
		if err := uc.repo.UpdateRefreshToken(user.GUID, refreshToken, publicKey); err != nil {
			return "", "", "", err
		}
		return accessToken, refreshToken, publicKey, err
	})

	encryptedRefreshToken := utils.EncodeBase64(refreshToken)
	return accessToken, encryptedRefreshToken, publicKey, nil
}

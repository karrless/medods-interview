// Package controllers provides tools for working with http requests
package controllers

import (
	"context"
	"medods-jwt/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthService is a service for users
type AuthService interface {
	GetAccessToken(guid, ip string) (accessToken, refreshToken, publicKey string, err error)
	RefreshToken(refresh, public, ip string) (accessToken, refreshToken, publicKey string, err error)
}

// UserAuthControllerController is a controller for users
type AuthController struct {
	ctx     *context.Context
	service AuthService
}

// NewAuthController returns a new user controller
func NewAuthController(ctx *context.Context, service AuthService) *AuthController {
	return &AuthController{
		ctx:     ctx,
		service: service,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// GetAccessTokenRequest is a request for get access token
type GetAccessTokenRequest struct {
	GUID string `json:"guid"`
}

type GetAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	PublicKey    string `json:"public_key"`
}

// @Summary Get access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body GetAccessTokenRequest true "Get access token request"
// @Success 200 {object} GetAccessTokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /auth/token [post]
func (uc *AuthController) GetAccessToken(c *gin.Context) {
	var req GetAccessTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, publicKey, err := uc.service.GetAccessToken(req.GUID, c.ClientIP())

	if err != nil {
		var status int
		switch err {
		case errors.ErrInvalidGUID:
			status = http.StatusBadRequest
		case errors.ErrUserNotFound:
			status = http.StatusNotFound
		case errors.ErrUserAlreadyHasAccessToken:
			status = http.StatusConflict
		default:
			status = http.StatusInternalServerError
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	response := GetAccessTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		PublicKey:    publicKey,
	}
	c.JSON(http.StatusOK, response)
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
	PublicKey    string `json:"public_key"`
}
type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	PublicKey    string `json:"public_key"`
}

// @Summary Refresh access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh access token request"
// @Success 200 {object} RefreshResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/refresh [post]
func (uc *AuthController) Refresh(c *gin.Context) {
	var req RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	accessToken, refreshToken, publicKey, err := uc.service.RefreshToken(req.RefreshToken, req.PublicKey, c.ClientIP())
	if err != nil {
		var status int
		switch err {
		case errors.ErrInvalidRefreshToken:
			status = http.StatusBadRequest
		case errors.ErrUserNotFound:
			status = http.StatusUnauthorized
		default:
			status = http.StatusInternalServerError
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	response := RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		PublicKey:    publicKey,
	}
	c.JSON(http.StatusOK, response)
}

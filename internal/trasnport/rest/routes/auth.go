package routes

import (
	"context"
	"medods-jwt/internal/service"
	"medods-jwt/internal/trasnport/rest/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(ctx *context.Context, r *gin.RouterGroup, authService *service.AuthService) {
	authController := controllers.NewAuthController(ctx, authService)
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/token", authController.GetAccessToken)
		authGroup.POST("/refresh", authController.Refresh)
	}
}

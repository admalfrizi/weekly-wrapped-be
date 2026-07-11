package router

import (
	"github.com/admalfrizi/weekly-wrapped-be/internal/controller"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(rg *gin.RouterGroup, authController *controller.AuthController) {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/login", authController.Login)
	}
}
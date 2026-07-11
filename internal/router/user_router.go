package router

import (
	"github.com/admalfrizi/weekly-wrapped-be/internal/controller"
	"github.com/admalfrizi/weekly-wrapped-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(
	rg *gin.RouterGroup, 
	userController *controller.UserController,
) {
	userRoutes := rg.Group("/users")
	userRoutes.Use(middleware.JWTMiddleware())
	{
		userRoutes.GET("/me", userController.GetProfile)
		userRoutes.POST("/me", userController.UpdateProfile)
	}
}
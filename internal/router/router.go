package router

import (
	"net/http"

	"github.com/admalfrizi/weekly-wrapped-be/internal/controller"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authController *controller.AuthController, 
	userController *controller.UserController,
	activityController *controller.ActivityController,
	dashboardController *controller.DashboardController,
) *gin.Engine {
	r := gin.Default();

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		});

		SetupAuthRoutes(v1, authController)

		SetupUserRoutes(v1, userController)

		SetupActivityRoutes(v1, activityController)

		SetupDashboardRoutes(v1, dashboardController)
	}

	return r
}
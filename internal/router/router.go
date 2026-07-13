package router

import (
	"net/http"
	"time"

	"github.com/admalfrizi/weekly-wrapped-be/internal/controller"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authController *controller.AuthController, 
	userController *controller.UserController,
	activityController *controller.ActivityController,
	dashboardController *controller.DashboardController,
	recapController *controller.RecapController,
) *gin.Engine {
	r := gin.Default();

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Allow your Next.js frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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

		SetupRecapRoutes(v1, recapController)
	}

	return r
}
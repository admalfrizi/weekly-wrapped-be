package router

import (
	"github.com/admalfrizi/weekly-wrapped-be/internal/controller"
	"github.com/admalfrizi/weekly-wrapped-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupActivityRoutes(
	rg *gin.RouterGroup,
	activityController *controller.ActivityController,
) {
	activityRoutes := rg.Group("/activity")
	activityRoutes.Use(middleware.JWTMiddleware())
	{
		activityRoutes.GET("/", activityController.GetActivities)
		activityRoutes.POST("/", activityController.Create)
		activityRoutes.GET("/:id", activityController.GetByID)
	}
}
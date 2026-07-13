package router

import (
	"github.com/admalfrizi/weekly-wrapped-be/internal/controller"
	"github.com/admalfrizi/weekly-wrapped-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupDashboardRoutes(rg *gin.RouterGroup, dashboardController *controller.DashboardController) {
	dashboardGroup := rg.Group("/dashboard")
	
	dashboardGroup.Use(middleware.JWTMiddleware())
	{
		dashboardGroup.GET("/weekly", dashboardController.GetWeekly)
	}
}
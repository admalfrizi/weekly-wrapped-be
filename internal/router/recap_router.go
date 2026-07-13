package router

import (
	"github.com/admalfrizi/weekly-wrapped-be/internal/controller"
	"github.com/admalfrizi/weekly-wrapped-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRecapRoutes(rg *gin.RouterGroup, recapController *controller.RecapController) {
	recapGroup := rg.Group("/recaps")

	recapGroup.Use(middleware.JWTMiddleware())
	{
		recapGroup.POST("/generate", recapController.Generate)
	}
	
	rg.GET("/recaps/:slug", recapController.GetBySlug)
}
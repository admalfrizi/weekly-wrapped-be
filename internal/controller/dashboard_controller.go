package controller

import (
	"net/http"
	"time"

	"github.com/admalfrizi/weekly-wrapped-be/internal/response"
	"github.com/admalfrizi/weekly-wrapped-be/internal/service"
	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	service service.DashboardService
}

func NewDashboardController(s service.DashboardService) *DashboardController {
	return &DashboardController{service: s}
}

func (c *DashboardController) GetWeekly(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Unauthorized", nil))
		return
	}

	dateStr := ctx.Query("start_date")
	startDate := time.Now()
	
	if dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Invalid start_date format. Use YYYY-MM-DD", nil))
			return
		}
		startDate = parsedDate
	}

	dashboardData, err := c.service.GetWeeklyDashboard(ctx.Request.Context(), userID, startDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to generate dashboard", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Weekly dashboard generated successfully", dashboardData))
}
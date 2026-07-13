package controller

import (
	"net/http"
	"time"

	"github.com/admalfrizi/weekly-wrapped-be/internal/response"
	"github.com/admalfrizi/weekly-wrapped-be/internal/service"
	"github.com/gin-gonic/gin"
)

type RecapController struct {
	service service.RecapService
}

func NewRecapController(s service.RecapService) *RecapController {
	return &RecapController{service: s}
}

// POST /api/v1/recaps/generate?start_date=2026-07-13
func (c *RecapController) Generate(ctx *gin.Context) {
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

	recap, err := c.service.GenerateRecap(ctx.Request.Context(), userID, startDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "Failed to generate recap", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Recap generated successfully", recap))
}

// GET /api/v1/recaps/:slug
func (c *RecapController) GetBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if slug == "" {
		ctx.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Slug is required", nil))
		return
	}

	recap, err := c.service.GetRecapBySlug(ctx.Request.Context(), slug)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Recap not found", err.Error()))
		return
	}

	// We decode the JSON snapshot so it returns as a clean object in the API response,
	// rather than an escaped string.
	ctx.JSON(http.StatusOK, response.Success("Recap retrieved successfully", recap))
}
package controller

import (
	"net/http"
	"strconv"

	"github.com/admalfrizi/weekly-wrapped-be/internal/response"
	"github.com/gin-gonic/gin"
)

type ActivityController struct {
}

func (c *ActivityController) GetActivities(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	mockActivities := []map[string]interface{}{
		{"id": 1, "category": "Coding", "value": 120},
		{"id": 2, "category": "Reading", "value": 45},
	}
	totalItems := 25

	totalPages := (totalItems + limit - 1) / limit 

	meta := response.PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}


	ctx.JSON(http.StatusOK, response.SuccessWithPagination(
		"Activities retrieved successfully",
		mockActivities,
		meta,
	))
}
package controller

import (
	"net/http"
	"strconv"

	"github.com/admalfrizi/weekly-wrapped-be/internal/dto"
	"github.com/admalfrizi/weekly-wrapped-be/internal/response"
	"github.com/admalfrizi/weekly-wrapped-be/internal/service"
	"github.com/gin-gonic/gin"
)

type ActivityController struct {
	service service.ActivityService
}

func NewActivityController(s service.ActivityService) *ActivityController {
	return &ActivityController{service: s}
}

func (c *ActivityController) Create(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Unauthorized", nil))
		return
	}

	var req dto.CreateActivityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "Invalid input data", err.Error()))
		return
	}

	activity, err := c.service.Create(ctx.Request.Context(), userID, req)
	if err != nil {
		ctx.JSON(http.StatusConflict, response.Error(
			http.StatusConflict,
			"Error Detected !",
			err.Error(),
		))
		return
	}

	ctx.JSON(http.StatusCreated, response.Success("Activity successfully created", activity))
}

func (c *ActivityController) GetActivities(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Unauthorized", nil))
		return
	}

	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 { page = 1 }
	if limit < 1 { limit = 10 }

	activities, totalItems, err := c.service.List(ctx.Request.Context(), userID, page, limit)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "user not found", err.Error()))
			return
		}
	}

	safeResponses := response.MapToActivityResponseList(activities)

	totalPages := (totalItems + limit - 1) / limit
	meta := response.PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	ctx.JSON(
		http.StatusOK, 
		response.SuccessWithPagination(
			"Activities retrieved successfully", 
			safeResponses, 
			meta,
		),
	)
}

func (c *ActivityController) GetByID(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Unauthorized", nil))
		return
	}

	activityID := ctx.Param("id")

	activity, err := c.service.GetByID(ctx.Request.Context(), userID, activityID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "Error Detected", err.Error()))
		return
	}

	safeResponse := response.MapToActivityResponse(*activity)
	ctx.JSON(http.StatusOK, response.Success("Activity retrieved successfully", safeResponse))
}

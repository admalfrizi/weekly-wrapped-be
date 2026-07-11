package controller

import (
	"net/http"

	"github.com/admalfrizi/weekly-wrapped-be/internal/dto"
	"github.com/admalfrizi/weekly-wrapped-be/internal/response"
	"github.com/admalfrizi/weekly-wrapped-be/internal/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	service service.UserService
}

func NewUserController(svc service.UserService) *UserController {
	return &UserController{
		service: svc,
	}
}

func (ctrl *UserController) GetProfile(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "unauthorized", "missing user id"))
		return
	}

	profile, err := ctrl.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "user not found", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "internal server error", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("User data successfully retrieved", profile))
}

func (ctrl *UserController) UpdateProfile(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "unauthorized", "missing user id"))
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid request body", err.Error()))
		return
	}

	updatedProfile, err := ctrl.service.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, "user not found", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, "internal server error", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("profile updated successfully", updatedProfile))
}
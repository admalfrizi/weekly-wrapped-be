package controller

import (
	"net/http"

	"github.com/admalfrizi/weekly-wrapped-be/internal/dto"
	"github.com/admalfrizi/weekly-wrapped-be/internal/response"
	"github.com/admalfrizi/weekly-wrapped-be/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service service.AuthService
}

func NewAuthController(s service.AuthService) *AuthController {
	return &AuthController{service: s}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Error(
			http.StatusBadRequest,
			"Invalid input data",
			err.Error(),
		))
		return
	}

	user, err := c.service.Register(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusConflict, response.Error(
			http.StatusConflict,
			"Registration failed",
			err.Error(),
		))
		return
	}

	safeResponse := response.MapToUserResponse(*user)

	ctx.JSON(http.StatusCreated, response.Success(
		"User successfully registered",
		safeResponse,
	))
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Error(
			http.StatusBadRequest,
			"Invalid input data",
			err.Error(),
		))
		return
	}

	user, accessToken, refreshToken,accessExp, refreshExp, err := c.service.Login(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.Error(
			http.StatusUnauthorized,
			"Login failed",
			err.Error(),
		))
		return
	}

	safeResponse := response.MapToUserResponse(*user)

	ctx.JSON(http.StatusOK, response.Success(
		"User successfully logged in",
		gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"access_token_expires_at":  accessExp,
			"refresh_token_expires_at": refreshExp,
			"user":          safeResponse,
		},
	))
}

func (c *AuthController) Refresh(ctx *gin.Context) {
	var req dto.RefreshRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Error(
			http.StatusBadRequest,
			"Invalid input data",
			err.Error(),
		))
		return
	}

	accessToken, refreshToken, accessExp, refreshExp, err := c.service.Refresh(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.Error(
			http.StatusUnauthorized,
			"Token refresh failed",
			err.Error(),
		))
		return
	}

	ctx.JSON(http.StatusOK, response.Success(
		"Tokens successfully refreshed",
		gin.H{
			"access_token":             accessToken,
			"refresh_token":            refreshToken,
			"access_token_expires_at":  accessExp,  // The exact timestamp for FE
			"refresh_token_expires_at": refreshExp, // The exact timestamp for FE
		},
	))
}
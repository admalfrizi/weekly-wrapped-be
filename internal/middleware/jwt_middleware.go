package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/admalfrizi/weekly-wrapped-be/internal/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(
				http.StatusUnauthorized,
				"Authorization header is required",
				"Unauthorized",
			))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(
				http.StatusUnauthorized,
				"Invalid Authorization header format",
				"Must be Bearer <token>",
			))
			return
		}
		
		tokenString := parts[1]
		jwtSecret := os.Getenv("JWT_SECRET")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(
				http.StatusUnauthorized,
				"Invalid or expired token",
				err.Error(),
			))
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userIDInterface, exists := claims["sub"]; exists {
				if userIDStr, ok := userIDInterface.(string); ok {
					ctx.Set("userID", userIDStr)
				}
			}
		}

		ctx.Next()
	}
}
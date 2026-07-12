package service

import (
	"context"
	"errors"
	"time"

	"github.com/admalfrizi/weekly-wrapped-be/internal/dto"
	"github.com/admalfrizi/weekly-wrapped-be/internal/model"
	"github.com/admalfrizi/weekly-wrapped-be/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*model.User, error)
	Login(ctx context.Context, req dto.LoginRequest) (*model.User, string, string, int64, int64, error)
	Refresh(ctx context.Context, req dto.RefreshRequest) (string, string, int64, int64, error)
}

type authService struct {
	repo      repository.AuthRepository
	jwtSecret []byte
}

func NewAuthService(repo repository.AuthRepository, jwtSecret string) AuthService {
	return &authService{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*model.User, error) {
	_, err := s.repo.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.New("email already registered")
	}

	_, err = s.repo.FindByUsername(ctx, req.Username)
	if err == nil {
		return nil, errors.New("username already taken")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &model.User{
		Email:        req.Email,
		Username:     req.Username,
		Name:         req.Name,
		PasswordHash: string(hashedPassword),
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*model.User, string, string, int64, int64, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		if err.Error() == "user not found" {
			return nil, "", "", 0, 0, errors.New("invalid email or password")
		}
		return nil, "", "", 0, 0, err
	}
	
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, "", "", 0, 0, errors.New("invalid email or password")
	}

	accessToken, refreshToken, accessExp, refreshExp, err := s.generateTokens(user.ID.String())
	if err != nil {
		return nil, "", "", 0, 0, err
	}

	return user, accessToken, refreshToken, accessExp, refreshExp, nil
}

func (s *authService) Refresh(ctx context.Context, req dto.RefreshRequest) (string, string, int64, int64, error) {
	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", "", 0, 0, errors.New("invalid or expired refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", 0, 0, errors.New("invalid token claims")
	}

	userIDInterface, exists := claims["sub"]
	if !exists {
		return "", "", 0, 0, errors.New("missing subject in token")
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		return "", "", 0, 0, errors.New("invalid subject format")
	}

	return s.generateTokens(userIDStr)
}

func (s *authService) generateTokens(userID string) (string, string,int64, int64, error) {
	accessExp := time.Now().Add(time.Hour * 24 * 7).Unix()
	refreshExp := time.Now().Add(time.Hour * 24 * 7).Unix()

	accessClaims := jwt.MapClaims{
		"sub": userID,
		"exp": accessExp,
	}
	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := accessTokenObj.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", 0, 0, err
	}

	refreshClaims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refreshTokenObj.SignedString(s.jwtSecret)
	if err != nil {
		return "", "",0, 0, err
	}

	return accessToken, refreshToken, accessExp, refreshExp, nil
}
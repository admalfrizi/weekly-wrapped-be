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
	// Updated interface to return raw model and tokens
	Register(ctx context.Context, req dto.RegisterRequest) (*model.User, string, string, error)
	Login(ctx context.Context, req dto.LoginRequest) (*model.User, string, string, error)
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

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*model.User, string, string, error) {
	_, err := s.repo.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, "", "", errors.New("email already registered")
	}

	_, err = s.repo.FindByUsername(ctx, req.Username)
	if err == nil {
		return nil, "", "", errors.New("username already taken")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", "", errors.New("failed to hash password")
	}

	user := &model.User{
		Email:        req.Email,
		Username:     req.Username,
		Name:         req.Name,
		PasswordHash: string(hashedPassword),
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, "", "", err
	}

	accessToken, refreshToken, err := s.generateTokens(user.ID.String())
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*model.User, string, string, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		if err.Error() == "user not found" {
			return nil, "", "", errors.New("invalid email or password")
		}
		return nil, "", "", err
	}
	
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, "", "", errors.New("invalid email or password")
	}

	accessToken, refreshToken, err := s.generateTokens(user.ID.String())
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *authService) generateTokens(userID string) (string, string, error) {
	accessClaims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	}
	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := accessTokenObj.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", err
	}

	refreshClaims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refreshTokenObj.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
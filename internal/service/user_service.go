package service

import (
	"context"
	"errors"
	"github.com/admalfrizi/weekly-wrapped-be/internal/dto"
	"github.com/admalfrizi/weekly-wrapped-be/internal/repository"
)

type UserService interface {
	GetProfile(ctx context.Context, userID string) (*dto.UserProfileResponse, error)
	UpdateProfile(ctx context.Context, userID string, req dto.UpdateProfileRequest) (*dto.UserProfileResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetProfile(ctx context.Context, userID string) (*dto.UserProfileResponse, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &dto.UserProfileResponse{
		ID:       user.ID.String(),
		Email:    user.Email,
		Username: user.Username,
		Name:     user.Name,
	}, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID string, req dto.UpdateProfileRequest) (*dto.UserProfileResponse, error) {
	
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.Name = req.Name
	user.Username = req.Username

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, errors.New("failed to update profile")
	}

	return &dto.UserProfileResponse{
		ID:       user.ID.String(),
		Email:    user.Email,
		Username: user.Username,
		Name:     user.Name,
	}, nil
}
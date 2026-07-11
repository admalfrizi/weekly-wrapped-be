package repository

import (
	"context"
	"errors"

	"github.com/admalfrizi/weekly-wrapped-be/internal/model"
	"github.com/admalfrizi/weekly-wrapped-be/internal/query"
	"github.com/jackc/pgx/v5"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
}

type authRepository struct {
	*BaseRepository
}

func NewAuthRepository(base *BaseRepository) AuthRepository {
	return &authRepository{BaseRepository: base}
}

func (r *authRepository) CreateUser(ctx context.Context, user *model.User) error {
	err := r.db.QueryRow(ctx, query.InsertUser, 
		user.Email, 
		user.Username, 
		user.Name, 
		user.PasswordHash,
	).Scan(
		&user.ID, 
		&user.CreatedAt, 
		&user.UpdatedAt,
	)

	return err
}

func (r *authRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	err := r.db.QueryRow(ctx, query.GetUserByEmail, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Name,
		&user.ProfileImgURL,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *authRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User

	err := r.db.QueryRow(ctx, query.GetUserByUsername, username).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Name,
		&user.ProfileImgURL,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
package repository

import (
	"context"
	"errors"

	"github.com/admalfrizi/weekly-wrapped-be/internal/model"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
}

type userRepository struct {
	*BaseRepository
}

func NewUserRepository(base *BaseRepository) UserRepository {
	return &userRepository{
		BaseRepository: base,
	}
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT id, email, username, name, password_hash 
		FROM users 
		WHERE id = $1
	`

	user := &model.User{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Name,
		&user.PasswordHash,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users 
		SET name = $1, username = $2
		WHERE id = $3
	`

	cmdTag, err := r.db.Exec(ctx, query, user.Name, user.Username, user.ID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("no rows affected or user not found")
	}

	return nil
}

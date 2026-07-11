package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)


type DBTX interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type BaseRepository struct {
	db DBTX
}

func NewBaseRepository(db DBTX) *BaseRepository {
	return &BaseRepository{db: db}
}

func (b *BaseRepository) WithTx(tx pgx.Tx) *BaseRepository {
	return &BaseRepository{db: tx}
}
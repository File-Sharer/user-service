package repository

import (
	"context"

	"github.com/File-Sharer/user-service/internal/model"
	"github.com/jackc/pgx/v5"
)

type User interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByLogin(ctx context.Context, login string) (*model.User, error)
	ExistsByLogin(ctx context.Context, login string) bool
}

type Repository struct {
	User
}

func New(db *pgx.Conn) *Repository {
	return &Repository{
		User: NewUserRepo(db),
	}
}

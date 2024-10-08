package repository

import (
	"context"

	"github.com/File-Sharer/user-service/internal/model"
	"github.com/File-Sharer/user-service/internal/repository/postgres"
	redisrepo "github.com/File-Sharer/user-service/internal/repository/redis"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type User interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByLogin(ctx context.Context, login string) (*model.User, error)
	ExistsByLogin(ctx context.Context, login string) bool
}

type Repository struct {
	Postgres *postgres.Repository
	Redis    *redisrepo.Repository
}

func New(db *pgx.Conn, rdb *redis.Client) *Repository {
	return &Repository{
		Postgres: postgres.New(db),
		Redis: redisrepo.New(rdb),
	}
}

package redisrepo

import (
	"context"
	"time"

	"github.com/File-Sharer/user-service/internal/model"
	"github.com/redis/go-redis/v9"
)

type User interface {
	Create(ctx context.Context, key string, value []byte, expiry time.Duration) error
	Find(ctx context.Context, key string) (*model.User, error)
}

type Repository struct {
	User
}

func New(rdb *redis.Client) *Repository {
	return &Repository{
		User: NewUserRepo(rdb),
	}
}

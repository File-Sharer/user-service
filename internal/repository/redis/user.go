package redisrepo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/File-Sharer/user-service/internal/model"
	"github.com/redis/go-redis/v9"
)

type UserRepo struct {
	rdb *redis.Client
}

func NewUserRepo(rdb *redis.Client) *UserRepo {
	return &UserRepo{rdb: rdb}
}

func (r *UserRepo) Create(ctx context.Context, key string, value []byte, expiry time.Duration) error {
	err := r.rdb.Set(ctx, key, value, expiry).Err()
	return err
}

func (r *UserRepo) Find(ctx context.Context, key string) (*model.User, error) {
	user, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var userDB model.User
	if err := json.Unmarshal([]byte(user), &userDB); err != nil {
		return nil, err
	}

	return &userDB, nil
}

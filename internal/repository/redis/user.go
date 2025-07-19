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

func (r *UserRepo) Create(ctx context.Context, key string, value model.User, expiry time.Duration) error {
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.rdb.Set(ctx, key, valueJSON, expiry).Err()
}

func (r *UserRepo) Find(ctx context.Context, key string) (*model.User, error) {
	userCache, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := json.Unmarshal([]byte(userCache), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

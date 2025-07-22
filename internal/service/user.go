package service

import (
	"context"
	"time"

	"github.com/File-Sharer/user-service/internal/model"
	"github.com/File-Sharer/user-service/internal/repository"
	"github.com/File-Sharer/user-service/internal/repository/redisrepo"
	"github.com/redis/go-redis/v9"
)

type UserService struct {
	repo *repository.Repository
	rdb *redis.Client
}

func NewUserService(repo *repository.Repository, rdb *redis.Client) *UserService {
	return &UserService{
		repo: repo,
		rdb: rdb,
	}
}

func (s *UserService) FindByID(ctx context.Context, id string) (*model.User, error) {
	userCache, err := redisrepo.Get[model.User](s.rdb, ctx, UserPrefix(id))
	if err == nil {
		return userCache, nil
	}
	if err != redis.Nil {
		return nil, err
	}

	user, err := s.repo.Postgres.User.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := redisrepo.SetJSON(s.rdb, ctx, UserPrefix(id), *user, time.Hour * 48); err != nil {
		return nil, err
	}

	return user, nil
}

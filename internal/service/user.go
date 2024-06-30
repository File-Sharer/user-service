package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/File-Sharer/user-service/internal/model"
	"github.com/File-Sharer/user-service/internal/repository"
	"github.com/redis/go-redis/v9"
)

type UserService struct {
	repo *repository.Repository
}

func NewUserService(repo *repository.Repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) FindByID(ctx context.Context, id string) (*model.User, error) {
	user, err := s.repo.Redis.User.Find(ctx, userPrefix + id)
	if err == nil {
		fmt.Println("HELLO USER FROM REDIS")
		return user, nil
	}

	if err != redis.Nil {
		return nil, err
	}

	userDB, err := s.repo.Postgres.User.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	userJSON, err := json.Marshal(userDB)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Redis.User.Create(ctx, userPrefix + id, userJSON, time.Hour * 48); err != nil {
		return nil, err
	}

	fmt.Println("HELLO USER FROM POSTGRES!")
	return userDB, nil
}

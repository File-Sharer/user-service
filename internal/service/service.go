package service

import (
	"context"

	"github.com/File-Sharer/user-service/internal/model"
	"github.com/File-Sharer/user-service/internal/repository"
	"github.com/redis/go-redis/v9"
)

type Auth interface {
	SignUp(ctx context.Context, user *model.User) (string, error)
	SignIn(ctx context.Context, user *model.User) (string, error)
}

type User interface {
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByLogin(ctx context.Context, login string) (*model.User, error)
}

type Service struct {
	Auth
	User
}

func New(repo *repository.Repository, rdb *redis.Client) *Service {
	return &Service{
		Auth: NewAuthService(repo),
		User: NewUserService(repo, rdb),
	}
}

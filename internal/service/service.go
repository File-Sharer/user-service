package service

import (
	"context"

	pb "github.com/File-Sharer/user-service/hasher_pbs"
	"github.com/File-Sharer/user-service/internal/model"
	"github.com/File-Sharer/user-service/internal/rabbitmq"
	"github.com/File-Sharer/user-service/internal/repository"
	"github.com/redis/go-redis/v9"
)

type Auth interface {
	SignUp(ctx context.Context, user *model.User) (*model.User, *model.JWTPair, error)
	SignIn(ctx context.Context, user *model.User) (*model.User, *model.JWTPair, error)
}

type User interface {
	FindByID(ctx context.Context, id string) (*model.User, error)
}

type Service struct {
	Auth
	User
}

func New(repo *repository.Repository, rabbitmq *rabbitmq.MQConn, hasherClient pb.HasherClient, rdb *redis.Client) *Service {
	userService := NewUserService(repo, rdb)

	return &Service{
		Auth: NewAuthService(repo, rabbitmq, hasherClient, userService),
		User: userService,
	}
}

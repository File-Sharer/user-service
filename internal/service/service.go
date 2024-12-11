package service

import (
	"context"

	pb "github.com/File-Sharer/user-service/hasher_pbs"
	"github.com/File-Sharer/user-service/internal/model"
	"github.com/File-Sharer/user-service/internal/repository"
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

func New(repo *repository.Repository, hasherClient pb.HasherClient) *Service {
	userService := NewUserService(repo)

	return &Service{
		Auth: NewAuthService(repo, hasherClient, userService),
		User: userService,
	}
}

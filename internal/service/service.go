package service

import (
	"context"

	pb "github.com/File-Sharer/user-service/hasher_pbs"
	"github.com/File-Sharer/user-service/internal/model"
	"github.com/File-Sharer/user-service/internal/repository"
)

type Auth interface {
	SignUp(ctx context.Context, user *model.User) (string, error)
	SignIn(ctx context.Context, user *model.User) (string, error)
}

type User interface {
	FindByID(ctx context.Context, id string) (*model.User, error)
}

type Service struct {
	Auth
	User
}

func New(repo *repository.Repository, hasherClient pb.HasherClient) *Service {
	return &Service{
		Auth: NewAuthService(repo, hasherClient),
		User: NewUserService(repo),
	}
}

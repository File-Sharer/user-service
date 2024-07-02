package service

import (
	"context"
	"os"
	"strings"
	"time"

	pb "github.com/File-Sharer/user-service/hasher_pbs"
	"github.com/File-Sharer/user-service/internal/model"
	"github.com/File-Sharer/user-service/internal/repository"
	"github.com/File-Sharer/user-service/pkg/auth"
	"github.com/jackc/pgx/v5"
)

type AuthService struct {
	repo *repository.Repository
	hasher pb.HasherClient
}

func NewAuthService(repo *repository.Repository, hasherClient pb.HasherClient) *AuthService {
	return &AuthService{
		repo: repo,
		hasher: hasherClient,
	}
}

func (s *AuthService) SignUp(ctx context.Context, user *model.User) (*model.User, string, error) {
	user.Login = strings.TrimSpace(strings.ToLower(user.Login))

	if s.repo.Postgres.User.ExistsByLogin(ctx, user.Login) {
		return nil, "", errLoginAlreadyTaken
	}

	passwordHash, err := auth.HashPassword([]byte(strings.TrimSpace(user.Password)))
	if err != nil {
		return nil, "", err
	}
	user.Password = passwordHash

	res, err := s.hasher.NewUID(ctx, &pb.NewUIDReq{UserLogin: user.Login})
	if !res.Ok {
		return nil, "", err
	}
	user.ID = res.GetUid()
	user.Role = "USER"
	user.DateAdded = time.Now()

	if err := s.repo.Postgres.User.Create(ctx, user); err != nil {
		return nil, "", nil
	}

	token, err := auth.GenerateToken(user.ID, []byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, "", err
	}

	return user.DTO(), token, nil
}

func (s *AuthService) SignIn(ctx context.Context, user *model.User) (*model.User, string, error) {
	userDB, err := s.repo.Postgres.User.FindByLogin(ctx, strings.TrimSpace(strings.ToLower(user.Login)))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, "", errInvalidCredentials
		}
		return nil, "", err
	}

	if !auth.VerifyPassword([]byte(userDB.Password), []byte(user.Password)) {
		return nil, "", errInvalidCredentials
	}

	token, err := auth.GenerateToken(userDB.ID, []byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, "", err
	}

	return userDB.DTO(), token, nil
}

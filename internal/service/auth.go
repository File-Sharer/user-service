package service

import (
	"context"
	"os"
	"strings"

	"github.com/File-Sharer/user-service/internal/model"
	"github.com/File-Sharer/user-service/internal/repository"
	"github.com/File-Sharer/user-service/pkg/auth"
	"github.com/File-Sharer/user-service/pkg/hasher"
	"github.com/jackc/pgx/v5"
)

type AuthService struct {
	repo *repository.Repository
}

func NewAuthService(repo *repository.Repository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) SignUp(ctx context.Context, user *model.User) (string, error) {
	user.Login = strings.TrimSpace(strings.ToLower(user.Login))

	if s.repo.User.ExistsByLogin(ctx, user.Login) {
		return "", errLoginAlreadyTaken
	}

	passwordHash, err := auth.HashPassword([]byte(strings.TrimSpace(user.Password)))
	if err != nil {
		return "", err
	}
	user.Password = passwordHash

	userID, err := hasher.GenerateUserID(user.Login)
	if err != nil {
		return "", nil
	}
	user.ID = userID

	if err := s.repo.User.Create(ctx, user); err != nil {
		return "", nil
	}

	token, err := auth.GenerateToken(user.ID, []byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) SignIn(ctx context.Context, user *model.User) (string, error) {
	userDB, err := s.repo.User.FindByLogin(ctx, strings.TrimSpace(strings.ToLower(user.Login)))
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", errInvalidCredentials
		}
		return "", err
	}

	if !auth.VerifyPassword([]byte(userDB.Password), []byte(user.Password)) {
		return "", errInvalidCredentials
	}

	token, err := auth.GenerateToken(userDB.ID, []byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return token, nil
}

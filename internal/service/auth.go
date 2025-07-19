package service

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	pb "github.com/File-Sharer/user-service/hasher_pbs"
	"github.com/File-Sharer/user-service/internal/model"
	"github.com/File-Sharer/user-service/internal/rabbitmq"
	"github.com/File-Sharer/user-service/internal/repository"
	"github.com/File-Sharer/user-service/pkg/auth"
	"github.com/jackc/pgx/v5"
)

type AuthService struct {
	repo *repository.Repository
	rabbitmq *rabbitmq.MQConn
	hasher pb.HasherClient
	userService User
}

func NewAuthService(repo *repository.Repository, rabbitmq *rabbitmq.MQConn, hasherClient pb.HasherClient, userService User) *AuthService {
	return &AuthService{
		repo: repo,
		rabbitmq: rabbitmq,
		hasher: hasherClient,
		userService: userService,
	}
}

func (s *AuthService) SignUp(ctx context.Context, user *model.User) (*model.User, *model.JWTPair, error) {
	user.Login = strings.TrimSpace(strings.ToLower(user.Login))

	if s.repo.Postgres.User.ExistsByLogin(ctx, user.Login) {
		return nil, nil, errLoginAlreadyTaken
	}

	passwordHash, err := auth.HashPassword([]byte(strings.TrimSpace(user.Password)))
	if err != nil {
		return nil, nil, err
	}
	user.Password = passwordHash

	res, err := s.hasher.NewUID(ctx, &pb.NewUIDReq{UserLogin: user.Login})
	if !res.GetOk() {
		return nil, nil, err
	}
	user.ID = res.GetUid()
	user.Role = "USER"
	user.DateAdded = time.Now()

	if err := s.repo.Postgres.User.Create(ctx, user); err != nil {
		return nil, nil, err
	}

	userCreatedMQ, err := json.Marshal(userCreated{UserID: user.ID})
	if err != nil {
		return nil, nil, err
	}
	if err := s.rabbitmq.PublishExchange(rabbitmq.USERS_CREATE_EXCHANGE, userCreatedMQ); err != nil {
		return nil, nil, err
	}

	jwtPairRes, err := s.hasher.GenerateJWTPair(ctx, &pb.GenerateJWTPairReq{Secret: os.Getenv("HASHER_SECRET"), UserId: user.ID, Role: user.Role})
	if err != nil {
		return nil, nil, err
	}

	return user.DTO(), &model.JWTPair{
		AccessToken: jwtPairRes.GetAccessToken(),
		RefreshToken: jwtPairRes.GetRefreshToken(),
	}, nil
}

type userCreated struct {
	UserID string `json:"userId"`
}

func (s *AuthService) SignIn(ctx context.Context, user *model.User) (*model.User, *model.JWTPair, error) {
	userDB, err := s.repo.Postgres.User.FindByLogin(ctx, strings.TrimSpace(strings.ToLower(user.Login)))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil, errInvalidCredentials
		}
		return nil, nil, err
	}

	if !auth.VerifyPassword([]byte(userDB.Password), []byte(user.Password)) {
		return nil, nil, errInvalidCredentials
	}

	jwtPairRes, err := s.hasher.GenerateJWTPair(ctx, &pb.GenerateJWTPairReq{Secret: os.Getenv("HASHER_SECRET"), UserId: userDB.ID, Role: userDB.Role})
	if err != nil {
		return nil, nil, err
	}

	return userDB.DTO(), &model.JWTPair{
		AccessToken: jwtPairRes.GetAccessToken(),
		RefreshToken: jwtPairRes.GetRefreshToken(),
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*model.JWTPair, *model.User, error) {
	decodedJwt, err := s.hasher.DecodeJWT(ctx, &pb.DecodeJWTReq{Secret: os.Getenv("HASHER_SECRET"), Jwt: refreshToken})
	if err != nil || !decodedJwt.Ok {
		return nil, nil, err
	}

	jwtPair, err := s.hasher.GenerateJWTPair(ctx, &pb.GenerateJWTPairReq{
		Secret: os.Getenv("HASHER_SECRET"),
		UserId: decodedJwt.GetUserId(),
		Role: decodedJwt.GetRole(),
	})
	if err != nil || !jwtPair.Ok {
		return nil, nil, err
	}

	user, err := s.userService.FindByID(ctx, decodedJwt.GetUserId())
	if err != nil {
		return nil, nil, err
	}

	return &model.JWTPair{
		AccessToken: jwtPair.GetAccessToken(),
		RefreshToken: jwtPair.GetRefreshToken(),
	}, user, nil
}

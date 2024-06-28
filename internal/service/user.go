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
	rdb *redis.Client
}

func NewUserService(repo *repository.Repository, rdb *redis.Client) *UserService {
	return &UserService{
		repo: repo,
		rdb: rdb,
	}
}

func (s *UserService) FindByID(ctx context.Context, id string) (*model.User, error) {
	user, err := s.rdb.Get(ctx, userPrefix + id).Result()
	if err != nil {
		if err == redis.Nil {
			userDB, err := s.repo.User.FindByID(ctx, id)
			if err != nil {
				return nil, err
			}

			userJSON, err := json.Marshal(userDB)
			if err != nil {
				return nil, err
			}

			if err := s.rdb.Set(ctx, userPrefix + id, userJSON, time.Hour * 24).Err(); err != nil {
				return nil, err
			}

			fmt.Println("HELLO USER FROM POSTGRES!")
			return userDB, nil
		}

		return nil, err
	}

	var userDB model.User
	if err := json.Unmarshal([]byte(user), &userDB); err != nil {
		return nil, err
	}

	fmt.Println("HELLO USER FROM REDIS")
	return &userDB, nil
}

func (s *UserService) FindByLogin(ctx context.Context, login string) (*model.User, error) {
	user, err := s.rdb.Get(ctx, userLoginPrefix + login).Result()
	if err != nil {
		if err == redis.Nil {
			userDB, err := s.repo.User.FindByLogin(ctx, login)
			if err != nil {
				return nil, err
			}

			userJSON, err := json.Marshal(userDB)
			if err != nil {
				return nil, err
			}

			if err := s.rdb.Set(ctx, userLoginPrefix + login, userJSON, time.Hour * 24).Err(); err != nil {
				return nil, err
			}

			fmt.Println("HELLO USER LOGIN FROM POSTGRES")
			return userDB, nil
		}

		return nil, err
	}

	var userDB model.User
	if err := json.Unmarshal([]byte(user), &userDB); err != nil {
		return nil, err
	}

	fmt.Println("HELLO USER LOGIN FROM REDIS")
	return &userDB, nil
}

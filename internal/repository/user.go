package repository

import (
	"context"

	"github.com/File-Sharer/user-service/internal/model"
	"github.com/jackc/pgx/v5"
)

type UserRepo struct {
	db *pgx.Conn
}

func NewUserRepo(db *pgx.Conn) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	_, err := r.db.Exec(ctx, "insert into users(id, login, password) values($1, $2, $3)", user.ID, user.Login, user.Password)
	return err
}

func (r *UserRepo) FindByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := r.db.QueryRow(ctx, "select id, login, password, role, date_added from users where id = $1", id).Scan(&user.ID, &user.Login, &user.Password, &user.Role, &user.DateAdded); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) FindByLogin(ctx context.Context, login string) (*model.User, error) {
	var user model.User
	if err := r.db.QueryRow(ctx, "select id, login, password, role, date_added from users where login = $1", login).Scan(&user.ID, &user.Login, &user.Password, &user.Role, &user.DateAdded); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) ExistsByLogin(ctx context.Context, login string) bool {
	var exists bool
	r.db.QueryRow(ctx, "select exists(select 1 from users where login = $1)", login).Scan(&exists)
	return exists
}

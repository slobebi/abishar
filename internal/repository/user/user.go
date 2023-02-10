package user

import (
	"abishar/internal/config"
	"abishar/internal/pkg/redigo"
	"context"
	"database/sql"
	"log"

	enUser "abishar/internal/entity/user"

	"github.com/jmoiron/sqlx"
)

type (
	redis interface {
		Del(keys ...string) error
		Get(key string) *redigo.Result
		Keys(key string) *redigo.Result
		Setex(key string, expireTime int, value interface{}) error
	}
)

type Repository struct {
	database *sqlx.DB
	redis    redis
	jwt      config.JWT
}

func NewRepository(
	db *sqlx.DB,
	redis redis,
	jwt config.JWT,
) *Repository {
	return &Repository{
		database: db,
		redis:    redis,
		jwt:      jwt,
	}
}

func (r *Repository) InsertUser(ctx context.Context, form enUser.RegisterForm) (*enUser.User, error) {

	user := &enUser.User{}
	var id int64
	var username string
	var isAdmin bool

	err := r.database.QueryRowContext(ctx, `
    insert into users
      (username, password, salt)
    values ($1, $2, $3)
    returning id, username, is_admin  
    `, form.Username, form.Password, form.Salt).Scan(&id, &username, &isAdmin)
	if err != nil {
		log.Printf("[InsertUser] Failed to insert user. err: %v", err)
		return user, err
	}

	user.ID = id
	user.IsAdmin = isAdmin
	user.Username = username

	return user, nil
}

func (r *Repository) CheckUsername(ctx context.Context, username string) (int64, error) {

	var userID int64

	err := r.database.QueryRowContext(ctx, `
    select id from users
      where username = $1
  `, username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return userID, nil
		}

		log.Printf("[CheckUsername] failed to check username. Err: %v", err)
		return 0, err
	}

	return userID, nil
}

func (r *Repository) GetByUsername(ctx context.Context, username string) (*enUser.User, error) {
	user := &enUser.User{}
	err := r.database.Get(user, `
    select id, username, is_admin, password, salt, wallet
    from users
    where username = $1
  `, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, nil
		}

		log.Printf("[GetByUsername] failed to get user. err: %v", err)
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetUserWallet(ctx context.Context, username string) (*enUser.UserWallet, error) {
	user := &enUser.UserWallet{}
	err := r.database.Get(user, `
    select id, username, wallet
    from users
    where username = $1
  `, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, nil
		}

		log.Printf("[GetUserWallet] failed to get user. err: %v", err)
		return nil, err
	}

	return user, nil
}

func (r *Repository) AddWallet(ctx context.Context, amount int64, userID int64) error {
	_, err := r.database.ExecContext(ctx, `
    update users
      set wallet = wallet + $1
    where id = $2
  `, amount, userID)
	if err != nil {
		log.Printf("[AddWallet] failed to user waller. err: %v", err)
		return err
	}

	return nil
}

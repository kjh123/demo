package data

import (
	"context"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserQueryer interface {
	Find(ctx context.Context, id int) (*User, error)
	Select(ctx context.Context) ([]User, error)
}

type UserCommander interface {
	Create(ctx context.Context, user User) error
}

type User struct {
	ID        int
	Name      string
	CreatedAt time.Time `db:"created_at"`
}

type UserRepository struct {
	*sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Find(ctx context.Context, id int) (*User, error) {
	var user User
	if err := r.GetContext(ctx, &user, "SELECT `id`,`name`, `created_at` FROM `users` WHERE `id` = ? LIMIT 1", id); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Select(ctx context.Context) ([]User, error) {
	var users []User
	if err := r.SelectContext(ctx, &users, "SELECT * FROM `users`"); err != nil {
		slog.ErrorContext(ctx, "user.select", "err", err)
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) Create(ctx context.Context, user User) error {
	_, err := r.ExecContext(ctx, "INSERT INTO `users` (`name`, `created_at`) VALUES (?, ?)", user.Name, user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

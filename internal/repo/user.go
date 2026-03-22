package repo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*User, error)
	AddUser(ctx context.Context, id int64, username string) error
	AddAdmin(ctx context.Context, id int64, username string) error
	Remove(ctx context.Context, id int64) error
	IsAdmin(ctx context.Context, id int64) (bool, error)
}

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*User, error) {
	var user User
	err := r.pool.QueryRow(ctx, getUserByIDQuery, id).
		Scan(&user.ID, &user.Username, &user.IsAdmin, &user.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) AddUser(ctx context.Context, id int64, username string) error {
	_, err := r.pool.Exec(ctx, addUserQuery, id, username)
	return err
}

func (r *userRepository) AddAdmin(ctx context.Context, id int64, username string) error {
	_, err := r.pool.Exec(ctx, addAdminQuery, id, username)
	return err
}

func (r *userRepository) Remove(ctx context.Context, id int64) error {
	_, err := r.pool.Exec(ctx, removeUserQuery, id)
	return err
}

func (r *userRepository) IsAdmin(ctx context.Context, id int64) (bool, error) {
	var isAdmin bool
	err := r.pool.QueryRow(ctx, isAdminQuery, id).Scan(&isAdmin)

	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return isAdmin, nil
}

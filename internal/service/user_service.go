package service

import (
	"context"

	"github.com/Nazhgam/aqquiryq_bot/internal/repo"
)

type UserService interface {
	IsAllowed(ctx context.Context, telegramID int64) (bool, error)
	IsAdmin(ctx context.Context, telegramID int64) (bool, error)

	AddUser(ctx context.Context, id int64, username string) error
	AddAdmin(ctx context.Context, id int64, username string) error
	RemoveUser(ctx context.Context, id int64) error
}

type userService struct {
	repo repo.UserRepository
}

func NewUserService(repo repo.UserRepository) UserService {
	return &userService{repo: repo}
}

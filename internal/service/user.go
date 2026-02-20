package service

import (
	"context"
	"fmt"
)

func (s *userService) IsAllowed(ctx context.Context, telegramID int64) (bool, error) {
	user, err := s.repo.GetByID(ctx, telegramID)
	if err != nil {
		return false, err
	}

	return user != nil, nil
}

func (s *userService) IsAdmin(ctx context.Context, telegramID int64) (bool, error) {
	return s.repo.IsAdmin(ctx, telegramID)
}

func (s *userService) AddUser(ctx context.Context, id int64, username string) error {
	if id <= 0 {
		return fmt.Errorf("invalid user id")
	}
	return s.repo.AddUser(ctx, id, username)
}

func (s *userService) AddAdmin(ctx context.Context, id int64, username string) error {
	if id <= 0 {
		return fmt.Errorf("invalid user id")
	}
	return s.repo.AddAdmin(ctx, id, username)
}

func (s *userService) RemoveUser(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid user id")
	}
	return s.repo.Remove(ctx, id)
}

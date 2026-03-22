package service

import (
	"context"

	"github.com/Nazhgam/aqquiryq_bot/internal/repo"
)

type ContentService interface {
	GetClasses(ctx context.Context) ([]int, error)
	GetQuarters(ctx context.Context, class int) ([]int, error)
	GetContents(ctx context.Context, class, quarter int) ([]repo.Content, error)
	GetContentsByClass(ctx context.Context, class int) ([]repo.Content, error)
	GetContent(ctx context.Context, id int64) (*repo.Content, error)

	AddContent(ctx context.Context, title, url string, class, quarter, lessonNumber int) (int64, error)
	DeleteContent(ctx context.Context, id int64) error
}

type contentService struct {
	repo repo.ContentRepository
}

func NewContentService(repo repo.ContentRepository) ContentService {
	return &contentService{repo: repo}
}

package service

import (
	"context"
	"fmt"

	"github.com/Nazhgam/aqquiryq_bot/internal/repo"
)

func (s *contentService) GetClasses(ctx context.Context) ([]int, error) {
	return s.repo.GetAvailableClasses(ctx)
}

func (s *contentService) GetQuarters(ctx context.Context, class int) ([]int, error) {
	return s.repo.GetQuartersByClass(ctx, class)
}

func (s *contentService) GetContents(ctx context.Context, class, quarter int) ([]repo.Content, error) {
	return s.repo.GetByClassAndQuarter(ctx, class, quarter)
}

func (s *contentService) GetContentsByClass(ctx context.Context, class int) ([]repo.Content, error) {
	return s.repo.GetContentByClass(ctx, class)
}

func (s *contentService) GetContent(ctx context.Context, id int64) (*repo.Content, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *contentService) AddContent(
	ctx context.Context,
	title, url string,
	class, quarter, lessonNumber int,
) (int64, error) {

	if title == "" {
		return 0, fmt.Errorf("title cannot be empty")
	}

	if url == "" {
		return 0, fmt.Errorf("canva url cannot be empty")
	}

	if class <= 0 {
		return 0, fmt.Errorf("class must be greater than 0")
	}

	if quarter < 1 || quarter > 4 {
		return 0, fmt.Errorf("quarter must be between 1 and 4")
	}

	if lessonNumber < 1 {
		return 0, fmt.Errorf("lessonNumber must be greater than 0")
	}

	content := &repo.Content{
		Title:        title,
		CanvaURL:     url,
		Class:        class,
		Quarter:      quarter,
		LessonNumber: lessonNumber,
	}

	return s.repo.AddContent(ctx, content)
}

func (s *contentService) DeleteContent(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid content id")
	}

	return s.repo.DeleteContent(ctx, id)
}

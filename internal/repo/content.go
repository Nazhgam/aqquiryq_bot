package repo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContentRepository interface {
	GetAvailableClasses(ctx context.Context) ([]int, error)
	GetQuartersByClass(ctx context.Context, class int) ([]int, error)
	GetByClassAndQuarter(ctx context.Context, class, quarter int) ([]Content, error)
	GetContentByClass(ctx context.Context, class int) ([]Content, error)
	GetByID(ctx context.Context, id int64) (*Content, error)

	AddContent(ctx context.Context, content *Content) (int64, error)
	DeleteContent(ctx context.Context, id int64) error
}

type contentRepository struct {
	pool *pgxpool.Pool
}

func NewContentRepository(pool *pgxpool.Pool) ContentRepository {
	return &contentRepository{pool: pool}
}

func (r *contentRepository) GetAvailableClasses(ctx context.Context) ([]int, error) {
	rows, err := r.pool.Query(ctx, getAvailableClassesQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []int
	for rows.Next() {
		var c int
		if err := rows.Scan(&c); err != nil {
			return nil, err
		}
		classes = append(classes, c)
	}

	return classes, nil
}

func (r *contentRepository) GetQuartersByClass(ctx context.Context, class int) ([]int, error) {
	query := getQuartersByClassQuery

	rows, err := r.pool.Query(ctx, query, class)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quarters []int
	for rows.Next() {
		var q int
		if err := rows.Scan(&q); err != nil {
			return nil, err
		}
		quarters = append(quarters, q)
	}

	return quarters, nil
}

func (r *contentRepository) GetByClassAndQuarter(ctx context.Context, class, quarter int) ([]Content, error) {
	query := getByClassAndQuarterQuery

	rows, err := r.pool.Query(ctx, query, class, quarter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Content
	for rows.Next() {
		var c Content
		if err := rows.Scan(&c.ID, &c.Title, &c.CanvaURL, &c.Class, &c.Quarter, &c.LessonNumber); err != nil {
			return nil, err
		}
		result = append(result, c)
	}

	return result, nil
}

func (r *contentRepository) GetContentByClass(ctx context.Context, class int) ([]Content, error) {
	query := getByClassQuery

	rows, err := r.pool.Query(ctx, query, class)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Content
	for rows.Next() {
		var c Content
		if err := rows.Scan(&c.ID, &c.Title, &c.CanvaURL, &c.Class, &c.Quarter, &c.LessonNumber); err != nil {
			return nil, err
		}
		result = append(result, c)
	}

	return result, nil
}

func (r *contentRepository) GetByID(ctx context.Context, id int64) (*Content, error) {
	query := getContentByIDQuery

	var c Content
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&c.ID, &c.Title, &c.CanvaURL, &c.Class, &c.Quarter, &c.LessonNumber)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *contentRepository) AddContent(ctx context.Context, content *Content) (int64, error) {
	var id int64
	err := r.pool.QueryRow(
		ctx,
		addContentQuery,
		content.Title,
		content.CanvaURL,
		content.Class,
		content.Quarter,
		content.LessonNumber,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *contentRepository) DeleteContent(ctx context.Context, id int64) error {
	_, err := r.pool.Exec(ctx, deleteContentQuery, id)
	return err
}

package url

import (
	"context"
	"fmt"
	"time"

	"github.com/dtt4h/go-url-shortener/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type URLRepository interface {
	FindByID(ctx context.Context, id int64) (*model.URL, error)
	FindByCode(ctx context.Context, code string) (*model.URL, error)
	Save(ctx context.Context, url *model.URL) (*model.URL, error)
	Delete(ctx context.Context, id int64) error
}

type urlRepository struct {
	db *pgxpool.Pool
}

func NewURLRepository(db *pgxpool.Pool) URLRepository {
	return &urlRepository{db: db}
}

func (r *urlRepository) FindByID(ctx context.Context, id int64) (*model.URL, error) {
	query := `SELECT id, short_code, original_url, created_at, expires_at, click_count
				FROM urls WHERE id = $1`

	var url model.URL
	var expiresAt *int64

	err := r.db.QueryRow(ctx, query, id).Scan(
		&url.ID, &url.ShortCode, &url.OriginalURL,
		&url.CreatedAt, &expiresAt, &url.ClickCount,
	)

	if err != nil {
		return nil, fmt.Errorf("find url: %w", err)
	}

	if expiresAt != nil {
		ts := time.Unix(*expiresAt, 0)
		url.ExpiresAt = &ts
	}

	return &url, nil
}

func (r *urlRepository) FindByCode(ctx context.Context, code string) (*model.URL, error) {
	query := `SELECT id, short_code, original_url, created_at, expires_at, click_count
				FROM urls WHERE short_code = $1`

	var url model.URL
	var expiresAt *int64

	err := r.db.QueryRow(ctx, query, code).Scan(
		&url.ID, &url.ShortCode, &url.OriginalURL,
		&url.CreatedAt, &expiresAt, &url.ClickCount,
	)

	if err != nil {
		return nil, fmt.Errorf("find by code :%w", err)
	}

	if expiresAt != nil {
		ts := time.Unix(*expiresAt, 0)
		url.ExpiresAt = &ts
	}

	return &url, nil
}
func (r *urlRepository) Save(ctx context.Context, url *model.URL) (*model.URL, error) {
	query := `INSERT INTO urls(short_code, original_url, created_at, expires_at, click_count) 
				VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var expiresAt *int64

	err := r.db.QueryRow(ctx, query,
		url.ShortCode,
		url.OriginalURL,
		url.CreatedAt,
		expiresAt,
		url.ClickCount).Scan(&url.ID)

	if err != nil {
		return nil, fmt.Errorf("save by : %w", err)
	}

	if expiresAt != nil {
		ts := time.Unix(*expiresAt, 0)
		url.ExpiresAt = &ts
	}
	return url, nil
}

func (r *urlRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM urls WHERE id =$1`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deleted url: %w", err)
	}
	return nil
}

package url

import (
	"context"

	"github.com/dtt4h/go-url-shortener/internal/model"
)

type URLRepository interface {
	FindByID(ctx context.Context, id int64) (*model.URL, error)
	FindByCode(ctx context.Context, code string) (*model.URL, error)
	Save(ctx context.Context, url *model.URL) (*model.URL, error)
	Delete(ctx context.Context, id int64) error
}

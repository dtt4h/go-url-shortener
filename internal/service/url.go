package service

import (
	"context"

	"github.com/dtt4h/go-url-shortener/internal/model"
)

type URLService interface {
	CreateLongURL(ctx context.Context, originalURL string) (*model.URL, error)
	GetByCode(ctx context.Context, code string) (*model.URL, error)
}

package service

import (
	"context"
	"time"

	"github.com/dtt4h/go-url-shortener/internal/model"
	urlRepo "github.com/dtt4h/go-url-shortener/internal/repository/url"
)

type URLService interface {
	CreateShortURL(ctx context.Context, originalURL string) (*model.URL, error)
	GetByCode(ctx context.Context, code string) (*model.URL, error)
	DeleteByCode(ctx context.Context, code string) error
	IncrementClickCount(ctx context.Context, code string) error
}

type urlService struct {
	repo   urlRepo.URLRepository
	events EventService
}

func NewURLService(repo urlRepo.URLRepository, events EventService) URLService {
	return &urlService{repo: repo, events: events}
}

func (s *urlService) CreateShortURL(ctx context.Context, originalURL string) (*model.URL, error) {
	if err := Validate(originalURL); err != nil {
		return nil, err
	}

	existing, err := s.repo.FindByOriginalURL(ctx, originalURL)
	if err == nil && existing != nil {
		return existing, nil
	}

	code := GenerateShortCode(originalURL)

	url := &model.URL{
		ShortCode:   code,
		OriginalURL: originalURL,
		CreatedAt:   time.Now(),
		ClickCount:  0,
	}

	url, err = s.repo.Save(ctx, url)
	if err != nil {
		return nil, err
	}

	_ = s.events.PublishURLCreated(ctx, model.URLEvent{
		ShortCode:   url.ShortCode,
		OriginalURL: url.OriginalURL,
	})

	return url, nil
}

func (s *urlService) GetByCode(ctx context.Context, code string) (*model.URL, error) {
	url, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if url == nil {
		return nil, ErrNotFound
	}

	return url, nil
}

func (s *urlService) DeleteByCode(ctx context.Context, code string) error {
	url, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return err
	}

	if url == nil {
		return ErrNotFound
	}

	_ = s.events.PublishURLDeleted(ctx, model.URLEvent{
		ShortCode:   url.ShortCode,
		OriginalURL: url.OriginalURL,
	})

	return s.repo.Delete(ctx, url.ID)
}

func (s *urlService) IncrementClickCount(ctx context.Context, code string) error {
	url, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return err
	}

	if url == nil {
		return ErrNotFound
	}

	return s.repo.IncrementClickCount(ctx, url.ID)
}

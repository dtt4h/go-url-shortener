package service

import (
	"context"
	"fmt"
	"time"

	"github.com/dtt4h/go-url-shortener/internal/model"
	urlRepo "github.com/dtt4h/go-url-shortener/internal/repository/url"
)

type URLService interface {
	CreateShortURL(ctx context.Context, originalURL string) (*model.URL, error)
	GetByCode(ctx context.Context, code string) (*model.URL, error)
	//Save(ctx context.Context, id int64) (*model.URL, error)
	Delete(ctx context.Context, id int64) error
}

type urlService struct {
	repo   urlRepo.URLRepository
	events EventService
}

func NewURLService(repo urlRepo.URLRepository, events EventService) URLService {
	return &urlService{repo: repo, events: events}
}

func (s *urlService) CreateShortURL(ctx context.Context, originalURL string) (*model.URL, error) {
	// TODO: Проверить, существует ли уже такой originalURL в БД перед созданием нового
	// TODO: Добавить обработку случая, когда сгенерированный код уже существует (коллизия)

	if err := Validate(originalURL); err != nil {
		return nil, err
	}

	code := GenerateShortCode(originalURL)

	url := &model.URL{
		ShortCode:   code,
		OriginalURL: originalURL,
		CreatedAt:   time.Now(),
		ClickCount:  0,
	}

	url, err := s.repo.Save(ctx, url)
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
		return nil, fmt.Errorf("not found")
	}

	return url, nil
}

func (s *urlService) Delete(ctx context.Context, id int64) error {
	url, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if url == nil {
		return fmt.Errorf("not found")
	}

	_ = s.events.PublishURLDeleted(ctx, model.URLEvent{
		ShortCode:   url.ShortCode,
		OriginalURL: url.OriginalURL,
	})

	return s.repo.Delete(ctx, id)
}

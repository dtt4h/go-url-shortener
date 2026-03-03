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
	repo urlRepo.URLRepository
}

func NewURLService(repo urlRepo.URLRepository) URLService {
	return &urlService{repo: repo}
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

	return s.repo.Save(ctx, url)
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
	// TODO: Добавить логирование удаления (для аудита)

	return s.repo.Delete(ctx, id)
}

// URL SHORTENING BUSINESS LOGIC
// URL validation, short code generation, expiration handling
package service

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/dtt4h/go-url-shortener/internal/models"
	"github.com/dtt4h/go-url-shortener/internal/repository"
	"github.com/dtt4h/go-url-shortener/pkg/hash"
)

var (
	ErrInvalidURL        = errors.New("invalid URL")
	ErrShortCodeConflict = errors.New("short code already exists")
)

type shortenerService struct {
	repo repository.URLRepository
}

func NewShortenerService(repo repository.URLRepository) Shortener {
	return &shortenerService{repo: repo}
}

func (s *shortenerService) Shorten(originalURL string) (*models.ShortURL, error) {
	if !isValidURL(originalURL) {
		return nil, ErrInvalidURL
	}

	if existing, err := s.repo.FindByURL(originalURL); err == nil {
		return existing, nil
	}

	shortCode, err := s.generateUniqueShortCode()
	if err != nil {
		return nil, err
	}

	shortURL := &models.ShortURL{
		ID:          hash.GenerateID(),
		OriginalURL: originalURL,
		ShortCode:   shortCode,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.Create(shortURL); err != nil {
		return nil, err
	}

	return shortURL, nil
}

func (s *shortenerService) Redirect(shortCode string) (string, error) {
	shortURL, err := s.repo.FindByShortCode(shortCode)
	if err != nil {
		return "", err
	}
	return shortURL.OriginalURL, nil
}

func (s *shortenerService) generateUniqueShortCode() (string, error) {
	for i := 0; i < 5; i++ {
		shortCode := hash.GenerateShortCode(6)
		if _, err := s.repo.FindByShortCode(shortCode); err != nil {
			return shortCode, nil
		}
	}
	return "", ErrShortCodeConflict
}

func isValidURL(rawURL string) bool {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	u, err := url.Parse(rawURL)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// IN-MEMORY STORAGE IMPLEMENTATION
// Fast storage for development and testing (no persistence)
package repository

import (
	"errors"
	"sync"

	"github.com/dtt4h/go-url-shortener/internal/models"
)

var (
	ErrURLNotFound      = errors.New("URL not found")
	ErrURLAlreadyExists = errors.New("URL already exists")
)

type MemoryRepo struct {
	urls map[string]*models.ShortURL
	mu   sync.RWMutex
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		urls: make(map[string]*models.ShortURL),
	}
}

func (m *MemoryRepo) Create(url *models.ShortURL) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.urls[url.ShortCode]; exists {
		return ErrURLAlreadyExists
	}

	m.urls[url.ShortCode] = url
	return nil
}

func (m *MemoryRepo) FindByShortCode(shortCode string) (*models.ShortURL, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	url, exists := m.urls[shortCode]
	if !exists {
		return nil, ErrURLNotFound
	}

	return url, nil
}

func (m *MemoryRepo) FindByURL(originalURL string) (*models.ShortURL, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, url := range m.urls {
		if url.OriginalURL == originalURL {
			return url, nil
		}
	}

	return nil, ErrURLNotFound
}

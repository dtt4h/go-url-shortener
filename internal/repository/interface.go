// STORAGE INTERFACES
// Contracts for diff storage implementations (memory, postgres)
package repository

import "github.com/dtt4h/go-url-shortener/internal/models"

type URLRepository interface {
	Create(url *models.ShortURL) error
	FindByShortCode(shortCode string) (*models.ShortURL, error)
	FindByURL(originalURL string) (*models.ShortURL, error)
}

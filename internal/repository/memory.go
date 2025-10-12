// IN-MEMORY STORAGE IMPLEMENTATION
// Fast storage for development and testing (no persistance)
package repository

import (
	"github.com/dtt4h/go-url-shortener/internal/models"
)

type MemoryRepo struct {
	urls map[string]*models.ShortURL
}

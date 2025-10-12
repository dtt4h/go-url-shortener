// DATA MODELS AND STRUCTURES
// URL models, req/res DTOs, db entities
package models

import "time"

type ShortURL struct {
	ID          string
	OriginalURL string
	ShortCode   string
	CreatedAt   time.Time
}

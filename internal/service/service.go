package service

import "github.com/dtt4h/go-url-shortener/internal/models"

type Shortener interface {
	Shorten(url string) (*models.ShortURL, error)
	Redirect(shortCode string) (string, error)
}

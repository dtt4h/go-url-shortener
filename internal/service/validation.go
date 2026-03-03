package service

import (
	"net/url"
	"strings"
)

const maxURLLength = 2048

func Validate(originalURL string) error {
	if originalURL == "" {
		return ErrEmptyURL
	}

	if len(originalURL) > maxURLLength {
		return ErrURLTooLong
	}

	parsed, err := url.Parse(originalURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ErrInvalidURL
	}

	if !isAllowedScheme(parsed.Scheme) {
		return ErrUnsupportedScheme
	}

	return nil
}

func isAllowedScheme(scheme string) bool {
	scheme = strings.ToLower(scheme)
	return scheme == "http" || scheme == "https"
}

var (
	ErrEmptyURL          = &ValidationError{Code: "EMPTY_URL", Message: "URL cannot be empty"}
	ErrURLTooLong        = &ValidationError{Code: "URL_TOO_LONG", Message: "URL exceeds maximum length"}
	ErrInvalidURL        = &ValidationError{Code: "INVALID_URL", Message: "Invalid URL format"}
	ErrUnsupportedScheme = &ValidationError{Code: "UNSUPPORTED_SCHEME", Message: "Only http and https are allowed"}
)

type ValidationError struct {
	Code    string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

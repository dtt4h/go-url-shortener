package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/dtt4h/go-url-shortener/internal/logger"
	"github.com/dtt4h/go-url-shortener/internal/service"
	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	service service.URLService
	baseURL string
	log     *logger.Logger
}

func NewURLHandler(s service.URLService, baseURL string, log *logger.Logger) *URLHandler {
	return &URLHandler{service: s, baseURL: baseURL, log: log}
}

type CreateRequest struct {
	URL string `json:"url" binding:"required"`
}

type CreateResponse struct {
	ShortURL   string `json:"short_url"`
	ExpiresAt  int64  `json:"expires_at"`
	ClickCount int64  `json:"click_count"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *URLHandler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn(c.Request.Context(), "invalid request", "error", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: ErrInvalidRequest.Error()})
		return
	}

	originalURL := service.NormalizeURL(req.URL)

	urlEntity, err := h.service.CreateShortURL(c.Request.Context(), originalURL)
	if err != nil {
		if errors.Is(err, service.ErrInvalidURL) {
			h.log.Warn(c.Request.Context(), "invalid url", "error", err.Error())
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		h.log.Error(c.Request.Context(), "create short url failed", "error", err.Error())
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: ErrInternal.Error()})
		return
	}

	shortURL := h.baseURL + "/" + urlEntity.ShortCode

	var expiresAt int64
	if urlEntity.ExpiresAt != nil {
		expiresAt = urlEntity.ExpiresAt.Unix()
	}

	h.log.Info(c.Request.Context(), "url created",
		"short_url", shortURL,
		"original_url", originalURL,
	)

	c.JSON(http.StatusCreated, CreateResponse{
		ShortURL:   shortURL,
		ExpiresAt:  expiresAt,
		ClickCount: urlEntity.ClickCount,
	})
}

func (h *URLHandler) Get(c *gin.Context) {
	code := c.Param("code")

	url, err := h.service.GetByCode(c.Request.Context(), code)
	if err != nil {
		h.log.Warn(c.Request.Context(), "url not found", "code", code)
		c.JSON(http.StatusNotFound, ErrorResponse{Error: ErrNotFound.Error()})
		return
	}

	if url.ExpiresAt != nil && url.ExpiresAt.Before(time.Now()) {
		h.log.Warn(c.Request.Context(), "url expired", "code", code)
		c.JSON(http.StatusGone, ErrorResponse{Error: ErrLinkExpired.Error()})
		return
	}

	go func() {
		_ = h.service.IncrementClickCount(c.Request.Context(), code)
	}()

	h.log.Info(c.Request.Context(), "redirect",
		"code", code,
		"original_url", url.OriginalURL,
		"ip", c.ClientIP(),
	)

	c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
}

func (h *URLHandler) Delete(c *gin.Context) {
	code := c.Param("code")

	err := h.service.DeleteByCode(c.Request.Context(), code)
	if err != nil {
		h.log.Warn(c.Request.Context(), "delete failed", "code", code, "error", err.Error())
		c.JSON(http.StatusNotFound, ErrorResponse{Error: ErrNotFound.Error()})
		return
	}

	h.log.Info(c.Request.Context(), "url deleted", "code", code)
	c.Status(http.StatusNoContent)
}
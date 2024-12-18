package dtos

import "time"

type CreateShortUrlRequest struct {
	OriginalURL       string `json:"original_url" binding:"required"`
	UserID            string
	ExpirationSeconds *int `json:"expiration_seconds,omitempty"`
}

type CreateShortUrlResponse struct {
	ID          string    `json:"id"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	ExpiresAt   time.Time `json:"expires_at"`
}

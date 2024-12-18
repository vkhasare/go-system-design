package services

import (
	"fmt"
	"time"
	"url-shortening/dtos"
	"url-shortening/entities"

	"github.com/google/uuid"
)

func (s *shortURLService) CreateShortURL(req dtos.CreateShortUrlRequest) (*dtos.CreateShortUrlResponse, error) {
	// Generate a unique short URL token (for demonstration, using UUID)
	shortToken := uuid.New().String()

	expiresAt := time.Now().UTC()
	if req.ExpirationSeconds != nil {
		expiresAt = expiresAt.Add(time.Duration(*req.ExpirationSeconds) * time.Second)
	} else {
		// Default expiration - for example, 30 days if not specified
		expiresAt = expiresAt.Add(30 * 24 * time.Hour)
	}

	shortUrlEntity := &entities.ShortURL{
		OriginalURL: req.OriginalURL,
		ShortURL:    shortToken,
		UserID:      &req.UserID,
		QRCode:      nil, // QR code generation if needed
		ExpiresAt:   expiresAt,
		CreatedBy:   req.UserID,
		CreatedDate: time.Now().UTC(),
	}

	err := s.repo.CreateShortURL(shortUrlEntity)
	if err != nil {
		return nil, err
	}

	response := &dtos.CreateShortUrlResponse{
		ID:          fmt.Sprintf("%d", shortUrlEntity.ID),
		ShortURL:    shortUrlEntity.ShortURL,
		OriginalURL: shortUrlEntity.OriginalURL,
		ExpiresAt:   shortUrlEntity.ExpiresAt,
	}

	return response, nil
}

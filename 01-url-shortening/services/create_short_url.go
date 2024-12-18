package services

import (
	"fmt"
	"time"
	"url-shortening/dtos"
	"url-shortening/mapper"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *shortURLService) CreateShortURL(c *gin.Context, req dtos.CreateShortUrlRequest) (*dtos.CreateShortUrlResponse, error) {
	// Map DTO to entity
	shortUrlEntity := mapper.CreateShortUrlRequestMapper(c, req)

	// Compute short URL
	shortToken := uuid.New().String()

	// Compute expiration timestamp
	expiresAt := time.Now().UTC()
	if req.ExpirationSeconds != nil {
		expiresAt = expiresAt.Add(time.Duration(*req.ExpirationSeconds) * time.Second)
	} else {
		// Default expiration - for example, 30 days if not specified
		expiresAt = expiresAt.Add(30 * 24 * time.Hour)
	}

	// Set remaining fields in entity
	shortUrlEntity.ShortURL = shortToken
	shortUrlEntity.ExpiresAt = expiresAt

	// Save to DB
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

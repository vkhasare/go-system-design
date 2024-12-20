package services

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"
	"url-shortening/dtos"
	"url-shortening/mapper"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func isValidURL(u string) bool {
	parsed, err := url.ParseRequestURI(u)
	if err != nil {
		log.Println(err)
		return false
	}

	// Check if scheme and host are present.
	// Valid schemes often include http, https, etc.
	if parsed.Scheme == "" || parsed.Host == "" {
		log.Println("Scheme:", parsed.Scheme, "Host:", parsed.Host)
		return false
	}

	// Optionally, ensure scheme is one of the commonly used ones (http/https).
	// Comment this out if you allow all URL schemes.
	if !strings.EqualFold(parsed.Scheme, "http") && !strings.EqualFold(parsed.Scheme, "https") {
		log.Println("Scheme:", parsed.Scheme)
		return false
	}

	return true
}

func (s *shortURLService) CreateShortURL(c *gin.Context, req dtos.CreateShortUrlRequest) (*dtos.CreateShortUrlResponse, error) {
	// Validate the original URL
	if !isValidURL(req.OriginalURL) {
		return nil, errors.New("Failed URL validation")
	}

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

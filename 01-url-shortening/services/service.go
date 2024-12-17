package services

import (
	"fmt"
	"time"
	"url-shortening/dtos"
	"url-shortening/entities"
	"url-shortening/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ShortURLService interface {
	CreateShortURL(c *gin.Context, req dtos.CreateShortUrlRequest) (*dtos.CreateShortUrlResponse, error)
}

type shortURLService struct {
	repo repositories.ShortURLRepository
}

func NewShortURLService(repo repositories.ShortURLRepository) ShortURLService {
	return &shortURLService{repo: repo}
}

func (s *shortURLService) CreateShortURL(c *gin.Context, req dtos.CreateShortUrlRequest) (*dtos.CreateShortUrlResponse, error) {
	// Generate a unique short URL token (for demonstration, using UUID)
	shortToken := uuid.New().String()

	email, exists := c.Value("email").(string)
	if !exists {
		return nil, fmt.Errorf("email not found in context")
	}

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
		// In a real application, you'd derive UserID and CreatedBy from context (e.g., authenticated user)
		UserID:      &email,
		QRCode:      nil, // QR code generation if needed
		ExpiresAt:   expiresAt,
		CreatedBy:   email, // placeholder, replace with actual user
		CreatedDate: time.Now().UTC(),
	}

	err := s.repo.CreateShortURL(shortUrlEntity)
	if err != nil {
		return nil, err
	}

	response := &dtos.CreateShortUrlResponse{
		ID:          string(rune(shortUrlEntity.ID)),
		ShortURL:    shortUrlEntity.ShortURL,
		OriginalURL: shortUrlEntity.OriginalURL,
		ExpiresAt:   shortUrlEntity.ExpiresAt,
	}

	return response, nil
}

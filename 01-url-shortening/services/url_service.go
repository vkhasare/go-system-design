package services

import (
	"url-shortening/dtos"
	"url-shortening/repositories"

	"github.com/gin-gonic/gin"
)

type ShortURLService interface {
	CreateShortURL(c *gin.Context, req dtos.CreateShortUrlRequest) (*dtos.CreateShortUrlResponse, error)
	DeleteShortURLByID(c *gin.Context, id uint64) error
	GetOriginalURL(s string) (string, error)
	GetQRCode(id uint64, imgFormat string) ([]byte, error)
}

type shortURLService struct {
	repo repositories.ShortURLRepository
}

func NewShortURLService(repo repositories.ShortURLRepository) ShortURLService {
	return &shortURLService{repo: repo}
}

package services

import (
	"url-shortening/dtos"
	"url-shortening/repositories"
	"url-shortening/storageio"

	"github.com/gin-gonic/gin"
)

type ShortURLService interface {
	CreateShortURL(c *gin.Context, req dtos.CreateShortUrlRequest) (*dtos.CreateShortUrlResponse, error)
	DeleteShortURLByID(c *gin.Context, id uint64) error
	GetOriginalURL(s string) (string, error)
	GenerateQRCode(c *gin.Context, id uint64, imgFormat string) (string, []byte, error)
}

type shortURLService struct {
	repo           repositories.ShortURLRepository
	storageHandler storageio.FileStorageHandler
}

func NewShortURLService(repo repositories.ShortURLRepository, s storageio.FileStorageHandler) ShortURLService {
	return &shortURLService{repo: repo, storageHandler: s}
}

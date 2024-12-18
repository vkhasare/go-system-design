package services

import (
	"url-shortening/dtos"
	"url-shortening/repositories"
)

type ShortURLService interface {
	CreateShortURL(req dtos.CreateShortUrlRequest) (*dtos.CreateShortUrlResponse, error)
	DeleteShortURLByID(id uint64) (*dtos.DeleteUrlResponse, error)
}

type shortURLService struct {
	repo repositories.ShortURLRepository
}

func NewShortURLService(repo repositories.ShortURLRepository) ShortURLService {
	return &shortURLService{repo: repo}
}

package repositories

import (
	"url-shortening/entities"

	"gorm.io/gorm"
)

type ShortURLRepository interface {
	CreateShortURL(shortUrl *entities.ShortURL) error
	// You can add more methods like GetByShortURL, etc.
}

type shortURLRepository struct {
	db *gorm.DB
}

func NewShortURLRepository(db *gorm.DB) ShortURLRepository {
	return &shortURLRepository{db: db}
}

func (r *shortURLRepository) CreateShortURL(shortUrl *entities.ShortURL) error {
	return r.db.Create(shortUrl).Error
}

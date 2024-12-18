package repositories

import (
	"url-shortening/entities"

	"gorm.io/gorm"
)

type ShortURLRepository interface {
	CreateShortURL(shortUrl *entities.ShortURL) error
	DeleteByID(id uint64) (int64, error)
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

func (r *shortURLRepository) DeleteByID(id uint64) (int64, error) {
	result := r.db.Delete(&entities.ShortURL{}, id)
	if result.Error != nil {
		return -1, result.Error
	}

	if result.RowsAffected == 0 {
		// No record found and deleted
		return 0, gorm.ErrRecordNotFound
	}
	return result.RowsAffected, nil
}

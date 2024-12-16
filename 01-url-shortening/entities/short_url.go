package entities

import "time"

type ShortURL struct {
	ID               uint64     `gorm:"primaryKey;autoIncrement"`
	OriginalURL      string     `gorm:"type:text"`
	ShortURL         string     `gorm:"type:text;not null;uniqueIndex"`
	UserID           *string    `gorm:"type:varchar(255)"`
	QRCode           *string    `gorm:"type:text"`
	ExpiresAt        time.Time  `gorm:"not null"`
	CreatedBy        string     `gorm:"type:varchar(255);not null"`
	CreatedDate      time.Time  `gorm:"autoCreateTime"`
	LastModifiedBy   *string    `gorm:"type:varchar(255)"`
	LastModifiedDate *time.Time `gorm:"autoUpdateTime"`
}

package mapper

import (
	"time"
	"url-shortening/dtos"
	"url-shortening/entities"

	"github.com/gin-gonic/gin"
)

//takes dto converts to entity.

func CreateShortUrlRequestMapper(c *gin.Context, req dtos.CreateShortUrlRequest) *entities.ShortURL {
	var email string
	userID, ok := c.Get("email")

	if ok {
		email = userID.(string)
	}

	return &entities.ShortURL{
		OriginalURL: req.OriginalURL,
		UserID:      &email,
		CreatedBy:   email,
		CreatedDate: time.Now().UTC(),
	}
}

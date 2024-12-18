package services

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

func (s *shortURLService) DeleteShortURLByID(c *gin.Context, id uint64) error {
	num_rows, err := s.repo.DeleteByID(id)

	switch {
	case num_rows < 0:
		return fmt.Errorf("Failed to delete: %w", err)

	case num_rows == 0:
		return errors.New("Record not found")
	}

	return nil
}

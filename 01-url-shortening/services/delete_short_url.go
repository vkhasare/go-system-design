package services

import (
	"errors"
	"fmt"
	"url-shortening/dtos"
)

func (s *shortURLService) DeleteShortURLByID(id uint64) (*dtos.DeleteUrlResponse, error) {
	num_rows, err := s.repo.DeleteByID(id)

	switch {
	case num_rows < 0:
		return nil, fmt.Errorf("Failed to delete: %w", err)

	case num_rows == 0:
		return nil, errors.New("Record not found")

	default:
		return &dtos.DeleteUrlResponse{
			ID:      fmt.Sprintf("%d", id),
			Deleted: true,
		}, nil
	}
}
package services

import (
	"errors"
	"log"
	"time"
)

func (s *shortURLService) GetOriginalURL(origUrl string) (string, error) {
	log.Println("Lookup:", origUrl)
	su, err := s.repo.FindByShortURL(origUrl)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Check expiration
	if time.Now().UTC().After(su.ExpiresAt) {
		return "", errors.New("Link expired") // Treat expired as not found
	}

	if su.OriginalURL == "" {
		return "", errors.New("Original URL not found")
	}

	return su.OriginalURL, nil
}

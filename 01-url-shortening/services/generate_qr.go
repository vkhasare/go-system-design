package services

import (
	"log"

	"rsc.io/qr"
)

func (s *shortURLService) GetQRCode(id uint64) ([]byte, error) {
	oriUrl, err := s.GetOriginalURLByID(id)
	if err != nil {
		log.Println("id", id, oriUrl, "", "err:", err.Error())
		return nil, err
	}

	pngImage, err := qr.Encode(oriUrl, qr.L)
	if err != nil {
		log.Println("Encode error", err.Error())
		return nil, err
	}

	return pngImage.PNG(), nil
}

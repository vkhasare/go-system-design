package services

import (
	"bufio"
	"bytes"
	"log"

	"image/jpeg"
	"image/png"

	svg "github.com/ajstarks/svgo"
	"rsc.io/qr"
)

func (s *shortURLService) GetQRCode(id uint64, imgFormat string) ([]byte, error) {
	oriUrl, err := s.GetOriginalURLByID(id)
	if err != nil {
		log.Println("id", id, oriUrl, "", "err:", err.Error())
		return nil, err
	}

	qrImage, err := qr.Encode(oriUrl, qr.L)
	if err != nil {
		log.Println("Encode error", err.Error())
		return nil, err
	}

	pngImg := qrImage.PNG()
	switch imgFormat {
	case "png":
		return pngImg, nil
	case "jpeg":
		src, err := png.Decode(bytes.NewReader(pngImg))
		var b bytes.Buffer
		w := bufio.NewWriter(&b)
		err = jpeg.Encode(w, src, &jpeg.Options{Quality: 100})
		if err == nil {
			log.Default().Printf("b.Bytes(): %+v", b)
			return b.Bytes(), nil
		}
		return nil, err
	case "svg":
		var width, height int
		width, height = qrImage.Size, qrImage.Size
		var b bytes.Buffer
		w := bufio.NewWriter(&b)
		canvas := svg.New(w)
		canvas.Start(width, height)
		canvas.Circle(width/2, height/2, 100)
		canvas.Image(0, 0, width, height, "https://github.com/ajstarks/svgo/blob/master/gophercolor128x128.png?raw=true")
		canvas.End()
		log.Default().Printf("b.Bytes(): %+v", b)
		return b.Bytes(), nil
	}
	return qrImage.PNG(), nil
}

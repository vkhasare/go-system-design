package services

import (
	"bufio"
	"bytes"
	"log"

	"image/color"
	"image/jpeg"
	"image/png"

	"github.com/dennwc/gotrace"
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

	pngBuf := qrImage.PNG()
	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	switch imgFormat {
	case "png":
		return pngBuf, nil

	case "jpeg":
		src, err := png.Decode(bytes.NewReader(pngBuf))
		err = jpeg.Encode(w, src, &jpeg.Options{Quality: 100})
		if err != nil {
			log.Default().Printf("JPEG Encode error")
			return nil, err
		}
		err = w.Flush()
		if err != nil {
			log.Default().Printf("JPEG Flush error", b)
			return nil, err
		}
		return b.Bytes(), nil

	case "svg":
		src, err := png.Decode(bytes.NewReader(pngBuf))
		bm := gotrace.NewBitmapFromImage(src, func(x, y int, c color.Color) bool {
			r, g, b, _ := c.RGBA()
			return r+g+b > 128
		})
		paths, err := gotrace.Trace(bm, nil)
		if err != nil {
			log.Default().Printf("SVG Tracing error")
			return nil, err
		}

		gotrace.WriteSvg(w, src.Bounds(), paths, "")
		if err != nil {
			log.Default().Printf("SVG Write error")
			return nil, err
		}

		err = w.Flush()

		if err != nil {
			log.Default().Printf("SVG Flush error")
			return nil, err
		}

		return b.Bytes(), nil
	}
	return pngBuf, nil
}

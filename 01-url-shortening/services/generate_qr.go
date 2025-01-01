package services

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"

	"image/color"
	"image/jpeg"
	"image/png"

	"github.com/dennwc/gotrace"
	"github.com/gin-gonic/gin"
	"rsc.io/qr"
)

func getJPEGQRCode(pngBuf []byte, w io.Writer) error {
	src, err := png.Decode(bytes.NewReader(pngBuf))
	err = jpeg.Encode(w, src, &jpeg.Options{Quality: 100})
	if err != nil {
		log.Default().Printf("JPEG Encode error")
		return err
	}
	return nil
}

func getSVGQRCode(pngBuf []byte, w io.Writer) error {
	src, err := png.Decode(bytes.NewReader(pngBuf))
	bm := gotrace.NewBitmapFromImage(src, func(x, y int, c color.Color) bool {
		r, g, b, _ := c.RGBA()
		return r+g+b > 128
	})
	paths, err := gotrace.Trace(bm, nil)
	if err != nil {
		log.Default().Printf("SVG Tracing error")
		return err
	}

	err = gotrace.WriteSvg(w, src.Bounds(), paths, "")
	if err != nil {
		log.Default().Printf("SVG Write error")
		return err
	}
	return nil
}

func (s *shortURLService) uploadFile(c *gin.Context, b []byte, id uint64, imgFormat string) error {
	//Upload generated artifacts to storage server.
	bucketName := "qr-codes"
	objectName := fmt.Sprintf("%d.%s", id, imgFormat)
	if imgFormat == "" {
		objectName = fmt.Sprintf("%d.png", id)
	}
	err := s.storageHandler.UploadFile(c, bucketName, objectName, b, imgFormat)
	return err

}

func (s *shortURLService) GenerateQRCode(c *gin.Context, id uint64, imgFormat string) (string, []byte, error) {
	oriUrl, err := s.GetOriginalURLByID(id)
	if err != nil {
		log.Println("id", id, oriUrl, "", "err:", err.Error())
		return "", nil, err
	}

	qrImage, err := qr.Encode(oriUrl, qr.L)
	if err != nil {
		log.Println("Encode error", err.Error())
		return "", nil, err
	}

	pngBuf := qrImage.PNG()
	var contentType string
	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	switch imgFormat {
	default:
		fallthrough
	case "png":
		contentType = "image/png"
		w.Write(pngBuf)

	case "jpeg":
		contentType = "image/jpeg"
		err = getJPEGQRCode(pngBuf, w)
		if err != nil {
			return "", nil, err
		}

	case "svg":
		contentType = "image/svg+xml"
		err = getSVGQRCode(pngBuf, w)
		if err != nil {
			return "", nil, err
		}
	}

	err = w.Flush()

	if err != nil {
		log.Default().Printf(err.Error())
		return "", nil, errors.New("Buffer flush error")
	}

	err = s.uploadFile(c, b.Bytes(), id, imgFormat)
	if err != nil {
		log.Default().Printf("File upload failed. %v\n", err)
		return "", nil, errors.New("File upload failed.")
	}

	return contentType, b.Bytes(), nil
}

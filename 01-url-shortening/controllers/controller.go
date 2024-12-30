package controllers

import (
	"log"
	"net/http"
	"strconv"
	"url-shortening/dtos"
	"url-shortening/services"

	"github.com/gin-gonic/gin"
)

type URLController struct {
	service services.ShortURLService
}

func NewURLController(service services.ShortURLService) *URLController {
	return &URLController{
		service: service,
	}
}

// CreateShortURL handles POST /urls
func (ctrl *URLController) CreateShortURL(c *gin.Context) {
	var req dtos.CreateShortUrlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			ErrorCode:    "BAD_REQUEST",
			ErrorMessage: err.Error(),
		})
		return
	}

	resp, err := ctrl.service.CreateShortURL(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			ErrorCode:    "INTERNAL_ERROR",
			ErrorMessage: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// DeleteShortURL handles DELETE /urls/{id}
func (ctrl *URLController) DeleteShortURL(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			ErrorCode:    "INVALID_ID",
			ErrorMessage: "ID must be a numeric value",
		})
		return
	}

	log.Printf("Attempting to delete: %d", id)
	err = ctrl.service.DeleteShortURLByID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			ErrorCode:    "INTERNAL_ERROR",
			ErrorMessage: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// RedirectToOriginal handles GET /{short_url}
func (ctrl *URLController) RedirectToOriginal(c *gin.Context) {
	s := c.Param("shortUrl")

	originalURL, err := ctrl.service.GetOriginalURL(s)
	if err != nil {
		// TODO: Distinguish between 5xx and 404
		c.Status(http.StatusNotFound)
		return
	}

	// Redirect to the original URL
	c.Redirect(http.StatusFound, originalURL)
	log.Println("Redirected to", originalURL)
}

// GenerateQRCode handles POST /urls/{id}/qrcode
func (ctrl *URLController) GenerateQRCode(c *gin.Context) {
	log.Println("GenerateQRCode Handler")
	id := c.Param("id")

	var req dtos.GenerateQRCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("JSON bind failed.", req)
		c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			ErrorCode:    "BAD_REQUEST",
			ErrorMessage: err.Error(),
		})
		return
	}

	parsedId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		log.Println("Parsing error: ", err.Error())
		c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			ErrorCode:    "BAD_REQUEST",
			ErrorMessage: "Bad ID" + id,
		})
		return
	}

	qrCode, err := ctrl.service.GetQRCode(parsedId, req.ImageFormat)
	if err != nil {
		log.Println("ctrl.service.GetQRCode:", err.Error())
		c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			ErrorCode:    "BAD_REQUEST",
			ErrorMessage: err.Error(),
		})
		return
	}

	log.Default().Println("req.ImageFormat", req.ImageFormat)
	switch req.ImageFormat {
	case "png":
		c.Data(http.StatusOK, "image/png", qrCode)
	case "jpeg":
		c.Data(http.StatusOK, "image/jpeg", qrCode)
	case "svg":
		c.Data(http.StatusOK, "image/svg+xml", qrCode)
	default:
		c.Data(http.StatusOK, "image/png", qrCode)
	}

}

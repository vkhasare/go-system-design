package controllers

import (
	"net/http"
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

	resp, err := ctrl.service.CreateShortURL(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			ErrorCode:    "INTERNAL_ERROR",
			ErrorMessage: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

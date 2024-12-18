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
	userID, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			ErrorCode:    "BAD_REQUEST",
			ErrorMessage: "UserID missing",
		})
		return
	} else {
		if _, ok = userID.(string); !ok {
			c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
				ErrorCode:    "BAD_REQUEST",
				ErrorMessage: "Failed to extract UserID",
			})
			return
		}
		req.UserID = userID.(string)
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
	resp, err := ctrl.service.DeleteShortURLByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			ErrorCode:    "INTERNAL_ERROR",
			ErrorMessage: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

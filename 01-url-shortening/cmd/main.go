package main

import (
	"log"
	"net/http"
	"url-shortening/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	//Initialize middleware
	cfg := middleware.OIDCConfig{
		Issuer:   "http://localhost:8080/realms/url-shortner",
		ClientID: "account",
	}

	authMiddleware, err := middleware.AuthMiddleware(cfg)

	if err != nil {
		log.Fatalf("failed to create auth middleware: %v", err)
	}

	log.Println("Finished middleware initialization.")

	//Initialize gin router
	r := gin.Default()

	r.GET("/ping", authMiddleware, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"response": "pong"})
	})

	r.Run(":8053")
}

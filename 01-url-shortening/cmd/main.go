package main

import (
	"log"
	"os"
	"url-shortening/controllers"
	"url-shortening/middleware"
	"url-shortening/repositories"
	"url-shortening/services"
	"url-shortening/storageio"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	//Initialize middleware
	issuer := os.Getenv("KC_REALM_URL")
	if issuer == "" {
		issuer = "http://localhost:8080/realms/url-shortner" // Default value
	}

	cfg := middleware.OIDCConfig{
		Issuer:   issuer,
		ClientID: "account",
	}

	authMiddleware, err := middleware.AuthMiddleware(cfg)

	if err != nil {
		log.Fatalf("failed to create auth middleware: %v", err)
	}

	log.Println("Finished middleware initialization.")

	//Initialize db.
	dsn := "host=postgres user=keycloak password=password dbname=url-shortner port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	//Initialize file storage handler
	minioHandler, err := storageio.NewMinioFileStorageHandler(
		"minio-lb:9002", // endpoint
		"12TgVKAlZqiOWZrMEmXU",
		"FySQXC8JfyJV9TbHMNhlczHCIg4PLSpasZf2yJsD",
		false, // useSSL
	)
	if err != nil {
		log.Fatalf("Failed to create MinioFileStorageHandler: %v", err)
	}

	//Initialize 3-tier arch.
	repo := repositories.NewShortURLRepository(db)
	service := services.NewShortURLService(repo, minioHandler)
	controller := controllers.NewURLController(service)

	//Initialize gin router
	r := gin.Default()

	r.POST("/urls", authMiddleware, controller.CreateShortURL)
	r.DELETE("/urls/:id", authMiddleware, controller.DeleteShortURL)
	r.GET("/:shortUrl", controller.RedirectToOriginal)
	r.POST("/urls/:id/qrcode", authMiddleware, controller.GenerateQRCode)

	r.Run(":8053")
}

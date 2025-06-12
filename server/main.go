package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	"github.com/iamsuteerth/tx-qr-tool-backend/pkg/config"
	"github.com/iamsuteerth/tx-qr-tool-backend/pkg/controller"
	"github.com/iamsuteerth/tx-qr-tool-backend/pkg/middleware/cors"
	"github.com/iamsuteerth/tx-qr-tool-backend/pkg/middleware/request_id"
	"github.com/iamsuteerth/tx-qr-tool-backend/pkg/middleware/security"
	"github.com/iamsuteerth/tx-qr-tool-backend/pkg/repository"
	"github.com/iamsuteerth/tx-qr-tool-backend/pkg/service"
	"github.com/iamsuteerth/tx-qr-tool-backend/utils"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Warn().Msg("Warning: .env file not found")
	}

	utils.InitLogger()

	db := config.GetDBConnection()
	defer config.CloseDBConnection()

	registrationRepo := repository.NewRegistrationRepository(db)

	registrationService := service.NewRegistrationService(registrationRepo)

	registrationController := controller.NewRegistrationController(registrationService)

	router := gin.Default()

	router.Use(request_id.RequestIDMiddleware())
	router.Use(cors.SetupCORS())
	router.Use(security.APIKeyAuthMiddleware())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	router.POST("/register", registrationController.Register)
	router.GET("/download-csv", registrationController.DownloadCSV)

	router.NoRoute(func(c *gin.Context) {
		requestID := utils.GetRequestID(c)
		c.JSON(http.StatusNotFound, gin.H{
			"status":     "ERROR",
			"code":       "ROUTE_NOT_FOUND",
			"message":    "This route does not exist",
			"request_id": requestID,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Info().Str("port", port).Msg("Server starting")
	if err := router.Run(":" + port); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}

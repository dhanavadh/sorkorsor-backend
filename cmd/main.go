package main

import (
	"github.com/dhanavadh/sorkorsor-backend/config"
	"github.com/dhanavadh/sorkorsor-backend/database"
	"github.com/dhanavadh/sorkorsor-backend/handlers"
	"github.com/dhanavadh/sorkorsor-backend/repository"
	"github.com/dhanavadh/sorkorsor-backend/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		panic(err)
	}
	r2 := storage.NewR2Storage(cfg)

	messageRepo := repository.NewMessageRepository(db)

	messageHandler := handlers.NewMessageHandler(messageRepo)
	uploadHandler := handlers.NewUploadHandler(r2)

	gin.SetMode(cfg.GinMode)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	v1.GET("/messages", messageHandler.GetAllMessages)
	v1.GET("/message/:id", messageHandler.GetMessageById)
	v1.POST("/message", messageHandler.CreateMessage)
	v1.POST("/upload", uploadHandler.UploadFile)
	v1.POST("/upload/presign", uploadHandler.GetPresignedURLs)

	r.Run(":" + cfg.ServerPort)
}

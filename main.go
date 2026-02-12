package main

import (
	"log"

	"tcc-tech/queue-backend/config"
	"tcc-tech/queue-backend/database"
	"tcc-tech/queue-backend/handlers"
	"tcc-tech/queue-backend/routes"
	"tcc-tech/queue-backend/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()
	cfg := config.Get()

	db := database.InitDB(cfg.DSN)

	queueService := services.NewQueueService(db)
	queueHandler := handlers.NewQueueHandler(queueService)

	r := gin.Default()
	r.Use(cors.Default())

	routes.RegisterRoutes(r, queueHandler)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

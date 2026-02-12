package routes

import (
	"tcc-tech/queue-backend/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, qh *handlers.QueueHandler) {
	api := r.Group("/api")
	{
		api.GET("/health", qh.HandleHealthCheck)

		queue := api.Group("/queue")
		{
			queue.POST("/next", qh.HandleNextQueue)
			queue.GET("/current", qh.HandleGetCurrent)
			queue.POST("/clear", qh.HandleClearQueue)
		}
	}
}

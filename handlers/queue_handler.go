package handlers

import (
	"errors"
	"net/http"

	"tcc-tech/queue-backend/services"

	"github.com/gin-gonic/gin"
)

const dateTimeFormat = "02/01/2006 15:04" // DD/MM/YYYY HH:mm

type QueueHandler struct {
	service *services.QueueService
}

func NewQueueHandler(service *services.QueueService) *QueueHandler {
	return &QueueHandler{service: service}
}

// HandleNextQueue POST /api/queue/next — ออกบัตรคิวใหม่
func (h *QueueHandler) HandleNextQueue(c *gin.Context) {
	ticket, err := h.service.IssueNextTicket()
	if err != nil {
		if errors.Is(err, services.ErrQueueExhausted) {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to issue queue ticket",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"queue_number": ticket.QueueNumber,
		"issued_at":    ticket.IssuedAt.Format(dateTimeFormat),
	})
}

// HandleGetCurrent GET /api/queue/current — ดูคิวปัจจุบัน
func (h *QueueHandler) HandleGetCurrent(c *gin.Context) {
	queueNumber, issuedAt, err := h.service.GetCurrentQueue()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get current queue",
		})
		return
	}

	issuedAtStr := ""
	if issuedAt != nil {
		issuedAtStr = issuedAt.Format(dateTimeFormat)
	}

	c.JSON(http.StatusOK, gin.H{
		"queue_number": queueNumber,
		"issued_at":    issuedAtStr,
	})
}

// HandleClearQueue POST /api/queue/clear — ล้างคิว
func (h *QueueHandler) HandleClearQueue(c *gin.Context) {
	if err := h.service.ClearQueue(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to clear queue",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"queue_number": "00",
		"message":      "Queue has been cleared",
	})
}

// HandleHealthCheck GET /api/health — Health check
func (h *QueueHandler) HandleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

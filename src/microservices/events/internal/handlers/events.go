package handlers

import (
	"fmt"
	"net/http"
	"time"

	"events-service/internal/kafka"
	"events-service/internal/models"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	movieProducer   *kafka.Producer
	userProducer    *kafka.Producer
	paymentProducer *kafka.Producer
}

func NewEventHandler(movieProducer *kafka.Producer, userProducer *kafka.Producer, paymentProducer *kafka.Producer) *EventHandler {
	return &EventHandler{
		movieProducer:   movieProducer,
		userProducer:    userProducer,
		paymentProducer: paymentProducer,
	}
}

func (h *EventHandler) CreateMovieEvent(c *gin.Context) {
	var event models.MovieEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	key := fmt.Sprintf("%d", event.MovieId)
	partition, offset, err := h.movieProducer.SendMessage(key, event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to publish event: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.EventResponse{
		Status:    "success",
		Partition: partition,
		Offset:    offset,
		Event: models.Event{
			Id:        fmt.Sprintf("%d-%d", partition, offset),
			Type:      "movie",
			Timestamp: time.Now(),
			Payload:   make(map[string]interface{}),
		},
	})
}

func (h *EventHandler) CreateUserEvent(c *gin.Context) {
	var event models.UserEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	key := fmt.Sprintf("%d", event.UserId)
	partition, offset, err := h.userProducer.SendMessage(key, event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to publish event: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.EventResponse{
		Status:    "success",
		Partition: partition,
		Offset:    offset,
		Event: models.Event{
			Id:        fmt.Sprintf("%d-%d", partition, offset),
			Type:      "user",
			Timestamp: time.Now(),
			Payload:   make(map[string]interface{}),
		},
	})
}

func (h *EventHandler) CreatePaymentEvent(c *gin.Context) {
	var event models.PaymentEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	key := fmt.Sprintf("%d", event.PaymentId)
	partition, offset, err := h.paymentProducer.SendMessage(key, event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to publish event: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.EventResponse{
		Status:    "success",
		Partition: partition,
		Offset:    offset,
		Event: models.Event{
			Id:        fmt.Sprintf("%d-%d", partition, offset),
			Type:      "payment",
			Timestamp: time.Now(),
			Payload:   make(map[string]interface{}),
		},
	})
}

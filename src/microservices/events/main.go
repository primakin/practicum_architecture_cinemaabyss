package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"events-service/internal/handlers"
	"events-service/internal/kafka"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Printf("Events microservice start")

	broker := os.Getenv("KAFKA_BROKERS")
	if broker == "" {
		broker = "kafka:9092"
	}
	brokers := []string{broker}

	movieTopic := "movie-events"
	userTopic := "user-events"
	paymentTopic := "payment-events"

	consumers := []*kafka.Consumer{}

	for _, topic := range []string{movieTopic, userTopic, paymentTopic} {
		consumer, err := kafka.NewConsumer(brokers, topic)
		if err != nil {
			log.Fatalf("Failed to create consumer for %s: %v", topic, err)
		}
		defer consumer.Stop()
		consumers = append(consumers, consumer)

		if err := consumer.Start(); err != nil {
			log.Fatalf("Failed to start consumer for %s: %v", topic, err)
		}
	}

	movieProducer, err := kafka.NewProducer(brokers, movieTopic)
	if err != nil {
		log.Fatalf("Failed to create producer for %s: %v", movieTopic, err)
	}
	defer movieProducer.Close()

	userProducer, err := kafka.NewProducer(brokers, userTopic)
	if err != nil {
		log.Fatalf("Failed to create producer for %s: %v", userTopic, err)
	}
	defer userProducer.Close()

	paymentProducer, err := kafka.NewProducer(brokers, paymentTopic)
	if err != nil {
		log.Fatalf("Failed to create producer for %s: %v", paymentTopic, err)
	}
	defer paymentProducer.Close()

	eventHandler := handlers.NewEventHandler(movieProducer, userProducer, paymentProducer)

	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api")
	{
		events := api.Group("/events")
		{
			events.POST("/movie", eventHandler.CreateMovieEvent)
			events.POST("/user", eventHandler.CreateUserEvent)
			events.POST("/payment", eventHandler.CreatePaymentEvent)
			events.GET("/health", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"status": true,
				})
			})
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	go func() {
		if err := router.Run(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %s", port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}

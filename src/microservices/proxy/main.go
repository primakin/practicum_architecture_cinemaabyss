/*
 * CinemaAbyss API
 *
 * API спецификация для системы CinemaAbyss, включающая монолит и микросервисы.  Система CinemaAbyss представляет собой платформу для управления фильмами, пользователями, платежами и подписками. Архитектура системы включает монолитное приложение и выделенные микросервисы.
 *
 * API version: 1.0.0
 * Contact: support@cinemaabyss.com
 */

package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	openapi "github.com/GIT_USER_ID/GIT_REPO_ID/go"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func main() {
	log.Printf("Server started")

	port := getenv("PORT", "8000")
	monolithUrl := getenv("MONOLITH_URL", "http://monolith:8080")
	moviesServiceUrl := getenv("MOVIES_SERVICE_URL", "http://movies-service:8081")
	eventsServiceUrl := getenv("EVENTS_SERVICE_URL", "http://events-service:8082")
	gradualMigration, err := strconv.ParseBool(getenv("GRADUAL_MIGRATION", "true"))
	if err != nil {
		log.Fatal(err)
	}

	moviesMigrationPercent, err := strconv.Atoi(getenv("MOVIES_MIGRATION_PERCENT", "50"))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("API gateway port: %s", port)
	log.Printf("Monolith URL: %s", monolithUrl)
	log.Printf("Movies service URL: %s", moviesServiceUrl)
	log.Printf("Events service URL: %s", eventsServiceUrl)
	log.Printf("Gradual migration: %t", gradualMigration)
	log.Printf("Movies migration percent: %d", moviesMigrationPercent)

	EventsAPIService := openapi.NewEventsAPIService(eventsServiceUrl)
	EventsAPIController := openapi.NewEventsAPIController(EventsAPIService)

	HealthAPIService := openapi.NewHealthAPIService(moviesServiceUrl, eventsServiceUrl)
	HealthAPIController := openapi.NewHealthAPIController(HealthAPIService)

	MoviesAPIService := openapi.NewMoviesAPIService(monolithUrl, moviesServiceUrl, gradualMigration, moviesMigrationPercent)
	MoviesAPIController := openapi.NewMoviesAPIController(MoviesAPIService)

	PaymentsAPIService := openapi.NewPaymentsAPIService(monolithUrl)
	PaymentsAPIController := openapi.NewPaymentsAPIController(PaymentsAPIService)

	SubscriptionsAPIService := openapi.NewSubscriptionsAPIService(monolithUrl)
	SubscriptionsAPIController := openapi.NewSubscriptionsAPIController(SubscriptionsAPIService)

	UsersAPIService := openapi.NewUsersAPIService(monolithUrl)
	UsersAPIController := openapi.NewUsersAPIController(UsersAPIService)

	router := openapi.NewRouter(EventsAPIController, HealthAPIController, MoviesAPIController, PaymentsAPIController, SubscriptionsAPIController, UsersAPIController)

	log.Fatal(http.ListenAndServe(":"+port, router))
}

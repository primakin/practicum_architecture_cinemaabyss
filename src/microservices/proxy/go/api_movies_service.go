/*
 * CinemaAbyss API
 *
 * API спецификация для системы CinemaAbyss, включающая монолит и микросервисы.  Система CinemaAbyss представляет собой платформу для управления фильмами, пользователями, платежами и подписками. Архитектура системы включает монолитное приложение и выделенные микросервисы.
 *
 * API version: 1.0.0
 * Contact: support@cinemaabyss.com
 */

package openapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// MoviesAPIService is a service that implements the logic for the MoviesAPIServicer
// This service should implement the business logic for every endpoint for the MoviesAPI API.
// Include any external packages or services that will be required by this service.
type MoviesAPIService struct {
	monolithUrl            string
	moviesServiceUrl       string
	gradualMigration       bool
	moviesMigrationPercent int
	httpClient             *http.Client
	requestsSent           int
	requestsRedirected     int
}

// NewMoviesAPIService creates a default api service
func NewMoviesAPIService(monolithUrl string, moviesServiceUrl string, gradualMigration bool, moviesMigrationPercent int) *MoviesAPIService {
	return &MoviesAPIService{
		monolithUrl:            monolithUrl,
		moviesServiceUrl:       moviesServiceUrl,
		gradualMigration:       gradualMigration,
		moviesMigrationPercent: moviesMigrationPercent,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *MoviesAPIService) GetEndpointForNextRequest() string {
	counter := int(s.requestsSent*s.moviesMigrationPercent/100) + 1
	s.requestsSent++
	need_redirect := s.gradualMigration && (s.moviesMigrationPercent > 0) &&
		(counter > s.requestsRedirected || s.moviesMigrationPercent == 100)

	if need_redirect {
		s.requestsRedirected++
		return s.moviesServiceUrl
	}

	return s.monolithUrl
}

// GetAllMovies - Получение списка всех фильмов
func (s *MoviesAPIService) GetAllMovies(ctx context.Context) (ImplResponse, error) {
	url := fmt.Sprintf("%s/api/movies", s.GetEndpointForNextRequest())
	log.Printf("Request will be sent to %s", url)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return Response(http.StatusInternalServerError, Error{}), err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Response(resp.StatusCode, Error{}), fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	if _, err := decoder.Token(); err != nil {
		return Response(http.StatusInternalServerError, Error{}), err
	}

	var movies []Movie
	for decoder.More() {
		var movie Movie
		if err := decoder.Decode(&movie); err != nil {
			return Response(http.StatusInternalServerError, Error{}), err
		}
		movies = append(movies, movie)
	}

	return Response(http.StatusOK, movies), nil
}

// CreateMovie - Создание нового фильма
func (s *MoviesAPIService) CreateMovie(ctx context.Context, movieInput MovieInput) (ImplResponse, error) {
	url := fmt.Sprintf("%s/api/movies", s.GetEndpointForNextRequest())
	log.Printf("Request will be sent to %s", url)

	jsonData, err := json.Marshal(movieInput)
	if err != nil {
		return Response(http.StatusInternalServerError, Error{}), err
	}

	resp, err := s.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return Response(http.StatusInternalServerError, Error{}), err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return Response(resp.StatusCode, Error{}), fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	var movie Movie
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&movie); err != nil {
		return Response(http.StatusInternalServerError, Error{}), err
	}

	return Response(http.StatusCreated, movie), nil
}

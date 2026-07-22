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
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// HealthAPIService is a service that implements the logic for the HealthAPIServicer
// This service should implement the business logic for every endpoint for the HealthAPI API.
// Include any external packages or services that will be required by this service.
type HealthAPIService struct {
	moviesServiceUrl string
	eventsServiceUrl string
	httpClient       *http.Client
}

// NewHealthAPIService creates a default api service
func NewHealthAPIService(moviesServiceUrl string, eventsServiceUrl string) *HealthAPIService {
	return &HealthAPIService{
		moviesServiceUrl: moviesServiceUrl,
		eventsServiceUrl: eventsServiceUrl,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetProxyHealth - Проверка работоспособности API Gateway
func (s *HealthAPIService) GetProxyHealth(ctx context.Context) (ImplResponse, error) {
	return Response(http.StatusOK, "Strangler Fig Proxy is healthy"), nil
}

// GetMoviesServiceHealth - Проверка работоспособности микросервиса фильмов
func (s *HealthAPIService) GetMoviesServiceHealth(ctx context.Context) (ImplResponse, error) {
	url := fmt.Sprintf("%s/api/movies/health", s.moviesServiceUrl)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return Response(http.StatusInternalServerError, Error{}), err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Response(resp.StatusCode, Error{}), fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	var healthResponse GetMoviesServiceHealth200Response
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&healthResponse); err != nil {
		return Response(http.StatusInternalServerError, Error{}), err
	}

	return Response(http.StatusOK, healthResponse), nil
}

// GetEventsServiceHealth - Проверка работоспособности микросервиса событий
func (s *HealthAPIService) GetEventsServiceHealth(ctx context.Context) (ImplResponse, error) {
	url := fmt.Sprintf("%s/api/events/health", s.eventsServiceUrl)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return Response(http.StatusInternalServerError, Error{}), err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Response(resp.StatusCode, Error{}), fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	var healthResponse GetMoviesServiceHealth200Response
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&healthResponse); err != nil {
		return Response(http.StatusInternalServerError, Error{}), err
	}

	return Response(http.StatusOK, healthResponse), nil
}

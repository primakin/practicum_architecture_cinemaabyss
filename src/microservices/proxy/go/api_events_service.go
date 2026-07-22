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
	"net/http"
	"time"
)

// EventsAPIService is a service that implements the logic for the EventsAPIServicer
// This service should implement the business logic for every endpoint for the EventsAPI API.
// Include any external packages or services that will be required by this service.
type EventsAPIService struct {
	eventsServiceUrl string
	httpClient       *http.Client
}

// NewEventsAPIService creates a default api service
func NewEventsAPIService(eventsServiceUrl string) *EventsAPIService {
	return &EventsAPIService{
		eventsServiceUrl: eventsServiceUrl,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// CreateMovieEvent - Создание события фильма
func (s *EventsAPIService) CreateMovieEvent(ctx context.Context, movieEvent MovieEvent) (ImplResponse, error) {
	url := fmt.Sprintf("%s/api/events/movie", s.eventsServiceUrl)

	jsonData, err := json.Marshal(movieEvent)
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

	var eventResponse EventResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&eventResponse); err != nil {
		return Response(http.StatusInternalServerError, Error{}), err
	}

	return Response(http.StatusCreated, eventResponse), nil
}

// CreateUserEvent - Создание события пользователя
func (s *EventsAPIService) CreateUserEvent(ctx context.Context, userEvent UserEvent) (ImplResponse, error) {
	url := fmt.Sprintf("%s/api/events/user", s.eventsServiceUrl)

	jsonData, err := json.Marshal(userEvent)
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

	var eventResponse EventResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&eventResponse); err != nil {
		return Response(http.StatusInternalServerError, Error{}), err
	}

	return Response(http.StatusCreated, eventResponse), nil
}

// CreatePaymentEvent - Создание события платежа
func (s *EventsAPIService) CreatePaymentEvent(ctx context.Context, paymentEvent PaymentEvent) (ImplResponse, error) {
	url := fmt.Sprintf("%s/api/events/payment", s.eventsServiceUrl)

	jsonData, err := json.Marshal(paymentEvent)
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

	var eventResponse EventResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&eventResponse); err != nil {
		return Response(http.StatusInternalServerError, Error{}), err
	}

	return Response(http.StatusCreated, eventResponse), nil
}

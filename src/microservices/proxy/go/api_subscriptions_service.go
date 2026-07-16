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

// SubscriptionsAPIService is a service that implements the logic for the SubscriptionsAPIServicer
// This service should implement the business logic for every endpoint for the SubscriptionsAPI API.
// Include any external packages or services that will be required by this service.
type SubscriptionsAPIService struct {
	monolithUrl string
	httpClient  *http.Client
}

// NewSubscriptionsAPIService creates a default api service
func NewSubscriptionsAPIService(monolithUrl string) *SubscriptionsAPIService {
	return &SubscriptionsAPIService{
		monolithUrl: monolithUrl,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetAllSubscriptions - Получение списка всех подписок
func (s *SubscriptionsAPIService) GetAllSubscriptions(ctx context.Context, userId int32) (ImplResponse, error) {
	url := fmt.Sprintf("%s/api/subscriptions?user_id=%d", s.monolithUrl, userId)

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

	var subscriptions []Subscription
	for decoder.More() {
		var subscription Subscription
		if err := decoder.Decode(&subscription); err != nil {
			return Response(http.StatusInternalServerError, Error{}), err
		}
		subscriptions = append(subscriptions, subscription)
	}

	return Response(http.StatusOK, subscriptions), nil
}

// CreateSubscription - Создание новой подписки
func (s *SubscriptionsAPIService) CreateSubscription(ctx context.Context, subscriptionInput SubscriptionInput) (ImplResponse, error) {
	url := fmt.Sprintf("%s/api/subscriptions", s.monolithUrl)

	jsonData, err := json.Marshal(subscriptionInput)
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

	var subscription Subscription
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&subscription); err != nil {
		return Response(http.StatusInternalServerError, Error{}), err
	}

	return Response(http.StatusCreated, subscription), nil
}

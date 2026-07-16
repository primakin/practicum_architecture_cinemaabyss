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

// PaymentsAPIService is a service that implements the logic for the PaymentsAPIServicer
// This service should implement the business logic for every endpoint for the PaymentsAPI API.
// Include any external packages or services that will be required by this service.
type PaymentsAPIService struct {
	monolithUrl string
	httpClient  *http.Client
}

// NewPaymentsAPIService creates a default api service
func NewPaymentsAPIService(monolithUrl string) *PaymentsAPIService {
	return &PaymentsAPIService{
		monolithUrl: monolithUrl,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetAllPayments - Получение списка всех платежей
func (s *PaymentsAPIService) GetAllPayments(ctx context.Context, userId int32) (ImplResponse, error) {
	url := fmt.Sprintf("%s/api/payments?user_id=%d", s.monolithUrl, userId)

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

	var payments []Payment
	for decoder.More() {
		var payment Payment
		if err := decoder.Decode(&payment); err != nil {
			return Response(http.StatusInternalServerError, Error{}), err
		}
		payments = append(payments, payment)
	}

	return Response(http.StatusOK, payments), nil
}

// CreatePayment - Создание нового платежа
func (s *PaymentsAPIService) CreatePayment(ctx context.Context, paymentInput PaymentInput) (ImplResponse, error) {
	url := fmt.Sprintf("%s/api/payments", s.monolithUrl)

	jsonData, err := json.Marshal(paymentInput)
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

	var payment Payment
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&payment); err != nil {
		return Response(http.StatusInternalServerError, Error{}), err
	}

	return Response(http.StatusCreated, payment), nil
}

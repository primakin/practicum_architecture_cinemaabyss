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

// UsersAPIService is a service that implements the logic for the UsersAPIServicer
// This service should implement the business logic for every endpoint for the UsersAPI API.
// Include any external packages or services that will be required by this service.
type UsersAPIService struct {
	monolithUrl string
	httpClient  *http.Client
}

// NewUsersAPIService creates a default api service
func NewUsersAPIService(monolithUrl string) *UsersAPIService {
	return &UsersAPIService{
		monolithUrl: monolithUrl,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetAllUsers - Получение списка всех пользователей
func (s *UsersAPIService) GetAllUsers(ctx context.Context) (ImplResponse, error) {
	url := fmt.Sprintf("%s/api/users", s.monolithUrl)

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

	var users []User
	for decoder.More() {
		var user User
		if err := decoder.Decode(&user); err != nil {
			return Response(http.StatusInternalServerError, Error{}), err
		}
		users = append(users, user)
	}

	return Response(http.StatusOK, users), nil
}

// CreateUser - Создание нового пользователя
func (s *UsersAPIService) CreateUser(ctx context.Context, userInput UserInput) (ImplResponse, error) {
	url := fmt.Sprintf("%s/api/users", s.monolithUrl)

	jsonData, err := json.Marshal(userInput)
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

	var user User
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&user); err != nil {
		return Response(http.StatusInternalServerError, Error{}), err
	}

	return Response(http.StatusCreated, user), nil
}

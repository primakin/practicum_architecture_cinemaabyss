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
	"net/http"
	"strings"
)

// UsersAPIController binds http requests to an api service and writes the service results to the http response
type UsersAPIController struct {
	service      UsersAPIServicer
	errorHandler ErrorHandler
}

// UsersAPIOption for how the controller is set up.
type UsersAPIOption func(*UsersAPIController)

// WithUsersAPIErrorHandler inject ErrorHandler into controller
func WithUsersAPIErrorHandler(h ErrorHandler) UsersAPIOption {
	return func(c *UsersAPIController) {
		c.errorHandler = h
	}
}

// NewUsersAPIController creates a default api controller
func NewUsersAPIController(s UsersAPIServicer, opts ...UsersAPIOption) *UsersAPIController {
	controller := &UsersAPIController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the UsersAPIController
func (c *UsersAPIController) Routes() Routes {
	return Routes{
		"GetAllUsers": Route{
			"GetAllUsers",
			strings.ToUpper("Get"),
			"/api/users",
			c.GetAllUsers,
		},
		"CreateUser": Route{
			"CreateUser",
			strings.ToUpper("Post"),
			"/api/users",
			c.CreateUser,
		},
	}
}

// OrderedRoutes returns all the api routes in a deterministic order for the UsersAPIController
func (c *UsersAPIController) OrderedRoutes() []Route {
	return []Route{
		Route{
			"GetAllUsers",
			strings.ToUpper("Get"),
			"/api/users",
			c.GetAllUsers,
		},
		Route{
			"CreateUser",
			strings.ToUpper("Post"),
			"/api/users",
			c.CreateUser,
		},
	}
}

// GetAllUsers - Получение списка всех пользователей
func (c *UsersAPIController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.URL.Query().Has("id") {
		ctx = context.WithValue(ctx, "id", r.URL.Query().Get("id"))
	}
	result, err := c.service.GetAllUsers(ctx)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	_ = EncodeJSONResponse(result.Body, &result.Code, w)
}

// CreateUser - Создание нового пользователя
func (c *UsersAPIController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userInputParam UserInput
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&userInputParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertUserInputRequired(userInputParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertUserInputConstraints(userInputParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreateUser(r.Context(), userInputParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	_ = EncodeJSONResponse(result.Body, &result.Code, w)
}

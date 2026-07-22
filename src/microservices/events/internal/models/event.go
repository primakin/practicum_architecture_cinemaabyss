package models

import (
	"time"
)

type UserEvent struct {
	UserId    int32     `json:"user_id"`
	Username  string    `json:"username,omitempty"`
	Email     string    `json:"email,omitempty"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
}

type PaymentEvent struct {
	PaymentId  int32     `json:"payment_id"`
	UserId     int32     `json:"user_id"`
	Amount     float32   `json:"amount"`
	Status     string    `json:"status"`
	Timestamp  time.Time `json:"timestamp"`
	MethodType string    `json:"method_type,omitempty"`
}

type MovieEvent struct {
	MovieId     int32    `json:"movie_id"`
	Title       string   `json:"title"`
	Action      string   `json:"action"`
	UserId      int32    `json:"user_id,omitempty"`
	Rating      float32  `json:"rating,omitempty"`
	Genres      []string `json:"genres,omitempty"`
	Description string   `json:"description,omitempty"`
}

type Event struct {
	Id        string                 `json:"id"`
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload"`
}

type EventResponse struct {
	Status    string `json:"status"`
	Partition int32  `json:"partition"`
	Offset    int64  `json:"offset"`
	Event     Event  `json:"event"`
}

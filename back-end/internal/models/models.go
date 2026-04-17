package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Board struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

type Column struct {
	ID        string    `json:"id"`
	BoardID   string    `json:"board_id"`
	Title     string    `json:"title"`
	Position  int       `json:"position"`
	Cards     []Card    `json:"cards"`
	CreatedAt time.Time `json:"created_at"`
}

type Label struct {
	ID      string `json:"id"`
	BoardID string `json:"board_id"`
	Title   string `json:"title"`
	Color   string `json:"color"`
}

type Card struct {
	ID          string    `json:"id"`
	ColumnID    string    `json:"column_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Position    int       `json:"position"`
	Labels      []Label   `json:"labels"`
	CreatedAt   time.Time `json:"created_at"`
}

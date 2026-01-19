package routes

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(db *sql.DB) http.Handler {
	route := chi.NewRouter()

	return route
}

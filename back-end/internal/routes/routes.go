package routes

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"go-do-list/internal/controllers"
	appMiddleware "go-do-list/internal/middleware"
	"go-do-list/internal/services"
	ws "go-do-list/internal/websocket"
)

func SetupRoutes(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(appMiddleware.CORS)

	hub := ws.NewHub()
	go hub.Run()

	// Services
	authSvc := services.NewAuthService(db)
	// boardSvc := services.NewBoardService(db)
	// colSvc := services.NewColumnService(db)
	// cardSvc := services.NewCardService(db)

	// Controllers
	authCtrl := controllers.NewAuthController(authSvc)
	// boardCtrl := controllers.NewBoardController(boardSvc, hub)
	// colCtrl := controllers.NewColumnController(colSvc, boardSvc, hub)
	// cardCtrl := controllers.NewCardController(cardSvc, boardSvc, hub)

	// Auth routes (public)
	r.Post("/auth/register", authCtrl.Register)
	r.Post("/auth/login", authCtrl.Login)

	// WebSocket (auth via query param in production; keeping simple here)
	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWS(hub, w, r)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(appMiddleware.Auth)

		r.Get("/me", authCtrl.Me)

		// r.Get("/boards", boardCtrl.ListBoards)
		// r.Post("/boards", boardCtrl.CreateBoard)
		// r.Get("/boards/{id}", boardCtrl.GetBoardDetail)
		// r.Patch("/boards/{id}", boardCtrl.UpdateBoard)
		// r.Delete("/boards/{id}", boardCtrl.DeleteBoard)

		// r.Get("/boards/{id}/labels", boardCtrl.ListLabels)
		// r.Post("/boards/{id}/labels", boardCtrl.CreateLabel)
		// r.Delete("/boards/{id}/labels/{labelId}", boardCtrl.DeleteLabel)

		// r.Post("/boards/{id}/columns", colCtrl.CreateColumn)
		// r.Patch("/boards/{id}/columns/reorder", colCtrl.ReorderColumns)
		// r.Patch("/boards/{id}/columns/{colId}", colCtrl.UpdateColumn)
		// r.Delete("/boards/{id}/columns/{colId}", colCtrl.DeleteColumn)

		// r.Post("/boards/{id}/columns/{colId}/cards", cardCtrl.CreateCard)
		// r.Patch("/boards/{id}/cards/{cardId}", cardCtrl.UpdateCard)
		// r.Delete("/boards/{id}/cards/{cardId}", cardCtrl.DeleteCard)
		// r.Patch("/boards/{id}/cards/{cardId}/move", cardCtrl.MoveCard)

		// r.Post("/boards/{id}/cards/{cardId}/labels/{labelId}", cardCtrl.AddLabel)
		// r.Delete("/boards/{id}/cards/{cardId}/labels/{labelId}", cardCtrl.RemoveLabel)
	})

	return r
}

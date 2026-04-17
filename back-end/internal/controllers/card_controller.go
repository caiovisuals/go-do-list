package controllers

import (
	"encoding/json"
	"net/http"

	"go-do-list/internal/middleware"
	"go-do-list/internal/services"
	"go-do-list/internal/utils"
	ws "go-do-list/internal/websocket"

	"github.com/go-chi/chi/v5"
)

type CardController struct {
	svc      *services.CardService
	boardSvc *services.BoardService
	hub      *ws.Hub
}

func NewCardController(svc *services.CardService, boardSvc *services.BoardService, hub *ws.Hub) *CardController {
	return &CardController{svc: svc, boardSvc: boardSvc, hub: hub}
}

func (c *CardController) CreateCard(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	boardID := chi.URLParam(r, "id")
	colID := chi.URLParam(r, "colId")

	if _, err := c.boardSvc.GetBoard(boardID, userID); err != nil {
		utils.Error(w, http.StatusForbidden, "access denied")
		return
	}

	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	card, err := c.svc.CreateCard(colID, body.Title, body.Description)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	c.broadcast(boardID, userID)
	utils.JSON(w, http.StatusCreated, card)
}

func (c *CardController) UpdateCard(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	boardID := chi.URLParam(r, "id")
	cardID := chi.URLParam(r, "cardId")

	if _, err := c.boardSvc.GetBoard(boardID, userID); err != nil {
		utils.Error(w, http.StatusForbidden, "access denied")
		return
	}

	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	card, err := c.svc.UpdateCard(cardID, body.Title, body.Description)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	c.broadcast(boardID, userID)
	utils.JSON(w, http.StatusOK, card)
}

func (c *CardController) DeleteCard(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	boardID := chi.URLParam(r, "id")
	cardID := chi.URLParam(r, "cardId")

	if _, err := c.boardSvc.GetBoard(boardID, userID); err != nil {
		utils.Error(w, http.StatusForbidden, "access denied")
		return
	}

	if err := c.svc.DeleteCard(cardID); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	c.broadcast(boardID, userID)
	utils.JSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}

func (c *CardController) MoveCard(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	boardID := chi.URLParam(r, "id")
	cardID := chi.URLParam(r, "cardId")

	if _, err := c.boardSvc.GetBoard(boardID, userID); err != nil {
		utils.Error(w, http.StatusForbidden, "access denied")
		return
	}

	var body struct {
		ColumnID string `json:"column_id"`
		Position int    `json:"position"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	if err := c.svc.MoveCard(cardID, body.ColumnID, body.Position); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	c.broadcast(boardID, userID)
	utils.JSON(w, http.StatusOK, map[string]string{"message": "moved"})
}

func (c *CardController) AddLabel(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	boardID := chi.URLParam(r, "id")
	cardID := chi.URLParam(r, "cardId")
	labelID := chi.URLParam(r, "labelId")

	if _, err := c.boardSvc.GetBoard(boardID, userID); err != nil {
		utils.Error(w, http.StatusForbidden, "access denied")
		return
	}

	if err := c.svc.AddLabel(cardID, labelID); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	c.broadcast(boardID, userID)
	utils.JSON(w, http.StatusOK, map[string]string{"message": "label added"})
}

func (c *CardController) RemoveLabel(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	boardID := chi.URLParam(r, "id")
	cardID := chi.URLParam(r, "cardId")
	labelID := chi.URLParam(r, "labelId")

	if _, err := c.boardSvc.GetBoard(boardID, userID); err != nil {
		utils.Error(w, http.StatusForbidden, "access denied")
		return
	}

	if err := c.svc.RemoveLabel(cardID, labelID); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	c.broadcast(boardID, userID)
	utils.JSON(w, http.StatusOK, map[string]string{"message": "label removed"})
}

func (c *CardController) broadcast(boardID, userID string) {
	board, columns, labels, err := c.boardSvc.GetBoardDetail(boardID, userID)
	if err != nil {
		return
	}
	data, _ := json.Marshal(map[string]interface{}{
		"type":    "board_update",
		"board":   board,
		"columns": columns,
		"labels":  labels,
	})
	c.hub.Broadcast(boardID, data)
}

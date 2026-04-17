package controllers

import (
	"encoding/json"
	"net/http"

	"go-do-list/internal/middleware"
	"go-do-list/internal/services"
	"go-do-list/internal/utils"
)

type AuthController struct {
	svc *services.AuthService
}

func NewAuthController(svc *services.AuthService) *AuthController {
	return &AuthController{svc: svc}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var input services.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	user, err := c.svc.Register(input)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, user)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var input services.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	token, user, err := c.svc.Login(input)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

func (c *AuthController) Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	user, err := c.svc.GetByID(userID)
	if err != nil {
		utils.Error(w, http.StatusNotFound, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, user)
}

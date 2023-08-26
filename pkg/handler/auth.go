package handler

import (
	"net/http"

	"github.com/OrIX219/todo/pkg"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	var input todo.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		switch err.(type) {
		case *todo.ErrInvalidCredentials:
			newErrorResponse(c, http.StatusOK, err.Error())
			return
		default:
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, map[string]any{
		"token": token,
	})
}

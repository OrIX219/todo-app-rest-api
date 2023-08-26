package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authHeader = "Authorization"
	userCtx    = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "Empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid auth header")
		return
	}

	if headerParts[0] != "Bearer" {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid auth header")
		return
	}

	if headerParts[1] == "" {
		newErrorResponse(c, http.StatusUnauthorized, "Empty token")
		return
	}

	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, userId)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "User id not found")
		return 0, errors.New("User id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "Invalid user id type")
		return 0, errors.New("Invalid user id type")
	}

	return idInt, nil
}

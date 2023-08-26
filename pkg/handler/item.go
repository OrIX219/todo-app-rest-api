package handler

import (
	"net/http"
	"strconv"

	"github.com/OrIX219/todo/pkg"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid list id")
		return
	}

	var input todo.TodoItem
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	id, err := h.services.TodoItem.Create(userId, listId, input)
	if err != nil {
		var status int
		switch err.(type) {
		case *todo.ErrNoSuchList:
			status = http.StatusOK
		default:
			status = http.StatusInternalServerError
		}
		newErrorResponse(c, status, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"id": id,
	})
}

type getAllItemsResponse struct {
	Data []todo.TodoItem `json:"data"`
}

func (h *Handler) getAllItems(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid list id")
		return
	}

	items, err := h.services.TodoItem.GetAll(userId, listId)
	if err != nil {
		var status int
		switch err.(type) {
		case *todo.ErrNoSuchList:
			status = http.StatusOK
		default:
			status = http.StatusInternalServerError
		}
		newErrorResponse(c, status, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllItemsResponse{
		Data: items,
	})
}

func (h *Handler) getItemById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid item id")
		return
	}

	item, err := h.services.TodoItem.GetById(userId, itemId)
	if err != nil {
		var status int
		switch err.(type) {
		case *todo.ErrNoSuchItem:
			status = http.StatusOK
		default:
			status = http.StatusInternalServerError
		}
		newErrorResponse(c, status, err.Error())
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) updateItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid item id")
		return
	}

	var input todo.UpdateItemInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.TodoItem.Update(userId, id, input)
	if err != nil {
		var status int
		switch err.(type) {
		case *todo.ErrNoSuchItem:
			status = http.StatusOK
		case *todo.ErrInvalidUpdateItemInput:
			status = http.StatusBadRequest
		default:
			status = http.StatusInternalServerError
		}
		newErrorResponse(c, status, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
func (h *Handler) deleteItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid item id")
		return
	}

	err = h.services.TodoItem.Delete(userId, itemId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

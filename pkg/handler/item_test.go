package handler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	todo "github.com/OrIX219/todo/pkg"
	"github.com/OrIX219/todo/pkg/service"
	mock_service "github.com/OrIX219/todo/pkg/service/mock"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestHandler_createItem(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoItem, listId int,
		input todo.TodoItem)

	testTable := []struct {
		name             string
		inputListId      any
		inputBody        string
		inputItem        todo.TodoItem
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:        "OK",
			inputListId: 1,
			inputBody:   `{"title":"Test","description":"Description"}`,
			inputItem: todo.TodoItem{
				Title:       "Test",
				Description: "Description",
			},
			mockBehavior: func(s *mock_service.MockTodoItem, listId int,
				input todo.TodoItem) {
				s.EXPECT().Create(1, listId, input).Return(1, nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"id":1}`,
		},
		{
			name:        "Empty title",
			inputListId: 1,
			inputBody:   `{"title":""}`,
			inputItem: todo.TodoItem{
				Title: "",
			},
			mockBehavior: func(s *mock_service.MockTodoItem, listId int,
				input todo.TodoItem) {
			},
			expectedStatus:   400,
			expectedResponse: `{"message":"Invalid request body"}`,
		},
		{
			name:        "No title",
			inputListId: 1,
			inputBody:   `{"description":"Description"}`,
			inputItem: todo.TodoItem{
				Description: "Description",
			},
			mockBehavior: func(s *mock_service.MockTodoItem, listId int,
				input todo.TodoItem) {
			},
			expectedStatus:   400,
			expectedResponse: `{"message":"Invalid request body"}`,
		},
		{
			name:        "Invalid id",
			inputListId: "asd",
			inputBody:   `{"title":"Test","description":"Description"}`,
			inputItem: todo.TodoItem{
				Title:       "Test",
				Description: "Description",
			},
			mockBehavior: func(s *mock_service.MockTodoItem, listId int,
				input todo.TodoItem) {
			},
			expectedStatus:   400,
			expectedResponse: `{"message":"Invalid list id"}`,
		},
		{
			name:        "No list with such id",
			inputListId: 10,
			inputBody:   `{"title":"Test","description":"Description"}`,
			inputItem: todo.TodoItem{
				Title:       "Test",
				Description: "Description",
			},
			mockBehavior: func(s *mock_service.MockTodoItem, listId int,
				input todo.TodoItem) {
				s.EXPECT().Create(1, listId, input).Return(0,
					&todo.ErrNoSuchList{})
			},
			expectedStatus:   200,
			expectedResponse: `{"message":"No list with such id"}`,
		},
		{
			name:        "Service failure",
			inputListId: 1,
			inputBody:   `{"title":"Test","description":"Description"}`,
			inputItem: todo.TodoItem{
				Title:       "Test",
				Description: "Description",
			},
			mockBehavior: func(s *mock_service.MockTodoItem, listId int,
				input todo.TodoItem) {
				s.EXPECT().Create(1, listId, input).Return(0,
					errors.New("Service failure"))
			},
			expectedStatus:   500,
			expectedResponse: `{"message":"Service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			item := mock_service.NewMockTodoItem(c)
			if listId, ok := testCase.inputListId.(int); ok {
				testCase.mockBehavior(item, listId, testCase.inputItem)
			}

			services := &service.Service{TodoItem: item}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/api/lists/:id/items/", func(c *gin.Context) {
				c.Set(userCtx, 1)
			}, handler.createItem)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST",
				fmt.Sprintf("/api/lists/%v/items/", testCase.inputListId),
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.Equal(t, testCase.expectedResponse, w.Body.String())
		})
	}
}

func TestHandler_getItemById(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoItem, id int)

	testTable := []struct {
		name             string
		inputId          any
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:    "OK",
			inputId: 1,
			mockBehavior: func(s *mock_service.MockTodoItem, itemId int) {
				s.EXPECT().GetById(1, itemId).Return(todo.TodoItem{
					Id:          1,
					Title:       "Test",
					Description: "Description",
					Done:        true,
				}, nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"id":1,"title":"Test","description":"Description","done":true}`,
		},
		{
			name:             "Invalid id",
			inputId:          "asd",
			mockBehavior:     func(s *mock_service.MockTodoItem, itemId int) {},
			expectedStatus:   400,
			expectedResponse: `{"message":"Invalid item id"}`,
		},
		{
			name:    "No item with such id",
			inputId: 10,
			mockBehavior: func(s *mock_service.MockTodoItem, itemId int) {
				s.EXPECT().GetById(1, itemId).Return(todo.TodoItem{},
					&todo.ErrNoSuchItem{})
			},
			expectedStatus:   200,
			expectedResponse: `{"message":"No item with such id"}`,
		},
		{
			name:    "Service failure",
			inputId: 1,
			mockBehavior: func(s *mock_service.MockTodoItem, itemId int) {
				s.EXPECT().GetById(1, itemId).Return(todo.TodoItem{},
					errors.New("Service failure"))
			},
			expectedStatus:   500,
			expectedResponse: `{"message":"Service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			item := mock_service.NewMockTodoItem(c)
			if itemId, ok := testCase.inputId.(int); ok {
				testCase.mockBehavior(item, itemId)
			}

			services := &service.Service{TodoItem: item}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/api/items/:id", func(c *gin.Context) {
				c.Set(userCtx, 1)
			}, handler.getItemById)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET",
				fmt.Sprintf("/api/items/%v", testCase.inputId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.Equal(t, testCase.expectedResponse, w.Body.String())
		})
	}
}

func TestHandler_updateItem(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoItem, id int,
		input todo.UpdateItemInput)

	titleString := "Test"
	descriptionString := "Description"
	doneBool := true
	testTable := []struct {
		name             string
		inputId          any
		inputBody        string
		inputUpdate      todo.UpdateItemInput
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:      "OK",
			inputId:   1,
			inputBody: `{"title":"Test","description":"Description","done":true}`,
			inputUpdate: todo.UpdateItemInput{
				Title:       &titleString,
				Description: &descriptionString,
				Done:        &doneBool,
			},
			mockBehavior: func(s *mock_service.MockTodoItem, id int,
				input todo.UpdateItemInput) {
				s.EXPECT().Update(1, id, input).Return(nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"status":"ok"}`,
		},
		{
			name:      "Invalid id",
			inputId:   "asd",
			inputBody: `{"title":"Test","description":"Description","done":true}`,
			inputUpdate: todo.UpdateItemInput{
				Title:       &titleString,
				Description: &descriptionString,
				Done:        &doneBool,
			},
			mockBehavior: func(s *mock_service.MockTodoItem, id int,
				input todo.UpdateItemInput) {
			},
			expectedStatus:   400,
			expectedResponse: `{"message":"Invalid item id"}`,
		},
		{
			name:      "No item with such id",
			inputId:   10,
			inputBody: `{"title":"Test","description":"Description","done":true}`,
			inputUpdate: todo.UpdateItemInput{
				Title:       &titleString,
				Description: &descriptionString,
				Done:        &doneBool,
			},
			mockBehavior: func(s *mock_service.MockTodoItem, id int,
				input todo.UpdateItemInput) {
				s.EXPECT().Update(1, id, input).Return(&todo.ErrNoSuchItem{})
			},
			expectedStatus:   200,
			expectedResponse: `{"message":"No item with such id"}`,
		},
		{
			name:      "Invalid request body",
			inputId:   1,
			inputBody: `{}`,
			inputUpdate: todo.UpdateItemInput{
				Title:       nil,
				Description: nil,
				Done:        nil,
			},
			mockBehavior: func(s *mock_service.MockTodoItem, id int,
				input todo.UpdateItemInput) {
				s.EXPECT().Update(1, id, input).Return(
					&todo.ErrInvalidUpdateItemInput{})
			},
			expectedStatus:   400,
			expectedResponse: `{"message":"Invalid item update input"}`,
		},
		{
			name:      "Service failure",
			inputId:   1,
			inputBody: `{"title":"Test","description":"Description","done":true}`,
			inputUpdate: todo.UpdateItemInput{
				Title:       &titleString,
				Description: &descriptionString,
				Done:        &doneBool,
			},
			mockBehavior: func(s *mock_service.MockTodoItem, id int,
				input todo.UpdateItemInput) {
				s.EXPECT().Update(1, id, input).Return(
					errors.New("Service failure"))
			},
			expectedStatus:   500,
			expectedResponse: `{"message":"Service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			item := mock_service.NewMockTodoItem(c)
			if itemId, ok := testCase.inputId.(int); ok {
				testCase.mockBehavior(item, itemId, testCase.inputUpdate)
			}

			services := &service.Service{TodoItem: item}
			handler := NewHandler(services)

			r := gin.New()
			r.PUT("/api/items/:id", func(c *gin.Context) {
				c.Set(userCtx, 1)
			}, handler.updateItem)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT",
				fmt.Sprintf("/api/items/%v", testCase.inputId),
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.Equal(t, testCase.expectedResponse, w.Body.String())
		})
	}
}

func TestHandler_deleteItem(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoItem, id int)

	testTable := []struct {
		name             string
		inputId          any
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:    "OK",
			inputId: 1,
			mockBehavior: func(s *mock_service.MockTodoItem, id int) {
				s.EXPECT().Delete(1, id).Return(nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"status":"ok"}`,
		},
		{
			name:    "Invalid id",
			inputId: "asd",
			mockBehavior: func(s *mock_service.MockTodoItem, id int) {
			},
			expectedStatus:   400,
			expectedResponse: `{"message":"Invalid item id"}`,
		},
		{
			name:    "No item with such id",
			inputId: 10,
			mockBehavior: func(s *mock_service.MockTodoItem, id int) {
				s.EXPECT().Delete(1, id).Return(nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"status":"ok"}`,
		},
		{
			name:    "Service failure",
			inputId: 1,
			mockBehavior: func(s *mock_service.MockTodoItem, id int) {
				s.EXPECT().Delete(1, id).Return(errors.New("Service failure"))
			},
			expectedStatus:   500,
			expectedResponse: `{"message":"Service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			item := mock_service.NewMockTodoItem(c)
			if itemId, ok := testCase.inputId.(int); ok {
				testCase.mockBehavior(item, itemId)
			}

			services := &service.Service{TodoItem: item}
			handler := NewHandler(services)

			r := gin.New()
			r.DELETE("/api/items/:id", func(c *gin.Context) {
				c.Set(userCtx, 1)
			}, handler.deleteItem)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE",
				fmt.Sprintf("/api/items/%v", testCase.inputId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.Equal(t, testCase.expectedResponse, w.Body.String())
		})
	}
}

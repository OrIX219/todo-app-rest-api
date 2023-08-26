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

func TestHandler_createList(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoList, input todo.TodoList)

	testTable := []struct {
		name             string
		inputBody        string
		inputList        todo.TodoList
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:      "OK",
			inputBody: `{"title":"Test","description":"Description"}`,
			inputList: todo.TodoList{
				Title:       "Test",
				Description: "Description",
			},
			mockBehavior: func(s *mock_service.MockTodoList, input todo.TodoList) {
				s.EXPECT().Create(1, input).Return(1, nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"id":1}`,
		},
		{
			name:      "Empty title",
			inputBody: `{"title":""}`,
			inputList: todo.TodoList{
				Title: "",
			},
			mockBehavior:     func(s *mock_service.MockTodoList, input todo.TodoList) {},
			expectedStatus:   400,
			expectedResponse: `{"message":"Invalid request body"}`,
		},
		{
			name:      "No title",
			inputBody: `{"description":"Description"}`,
			inputList: todo.TodoList{
				Description: "Description",
			},
			mockBehavior:     func(s *mock_service.MockTodoList, input todo.TodoList) {},
			expectedStatus:   400,
			expectedResponse: `{"message":"Invalid request body"}`,
		},
		{
			name:      "Service failure",
			inputBody: `{"title":"Test","description":"Description"}`,
			inputList: todo.TodoList{
				Title:       "Test",
				Description: "Description",
			},
			mockBehavior: func(s *mock_service.MockTodoList, input todo.TodoList) {
				s.EXPECT().Create(1, input).Return(0, errors.New("Service failure"))
			},
			expectedStatus:   500,
			expectedResponse: `{"message":"Service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			list := mock_service.NewMockTodoList(c)
			testCase.mockBehavior(list, testCase.inputList)

			services := &service.Service{TodoList: list}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/api/lists/", func(c *gin.Context) {
				c.Set(userCtx, 1)
			}, handler.createList)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/lists/",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.Equal(t, testCase.expectedResponse, w.Body.String())
		})
	}
}

func TestHandler_getAllLists(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoList)

	testTable := []struct {
		name             string
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "OK",
			mockBehavior: func(s *mock_service.MockTodoList) {
				s.EXPECT().GetAll(1).Return([]todo.TodoList{
					{
						Id:          1,
						Title:       "Test",
						Description: "Description",
					},
					{
						Id:    2,
						Title: "Test2",
					},
				}, nil)
			},
			expectedStatus: 200,
			expectedResponse: "{\"data\":[" +
				"{\"id\":1,\"title\":\"Test\",\"description\":\"Description\"}," +
				"{\"id\":2,\"title\":\"Test2\",\"description\":\"\"}" +
				"]}",
		},
		{
			name: "No lists",
			mockBehavior: func(s *mock_service.MockTodoList) {
				s.EXPECT().GetAll(1).Return([]todo.TodoList{}, nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"data":[]}`,
		},
		{
			name: "Service failure",
			mockBehavior: func(s *mock_service.MockTodoList) {
				s.EXPECT().GetAll(1).Return([]todo.TodoList{},
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

			list := mock_service.NewMockTodoList(c)
			testCase.mockBehavior(list)

			services := &service.Service{TodoList: list}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/api/lists/", func(c *gin.Context) {
				c.Set(userCtx, 1)
			}, handler.getAllLists)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/lists/", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.Equal(t, testCase.expectedResponse, w.Body.String())
		})
	}
}

func TestHandler_getListById(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoList, id int)

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
			mockBehavior: func(s *mock_service.MockTodoList, id int) {
				s.EXPECT().GetById(1, id).Return(todo.TodoList{
					Id:          id,
					Title:       "Test",
					Description: "Description",
				}, nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"id":1,"title":"Test","description":"Description"}`,
		},
		{
			name:             "Invalid id",
			inputId:          "asd",
			mockBehavior:     func(s *mock_service.MockTodoList, id int) {},
			expectedStatus:   400,
			expectedResponse: `{"message":"Invalid list id"}`,
		},
		{
			name:    "No list with such id",
			inputId: 10,
			mockBehavior: func(s *mock_service.MockTodoList, id int) {
				s.EXPECT().GetById(1, id).Return(todo.TodoList{},
					&todo.ErrNoSuchList{})
			},
			expectedStatus:   200,
			expectedResponse: `{"message":"No list with such id"}`,
		},
		{
			name:    "Service failure",
			inputId: 1,
			mockBehavior: func(s *mock_service.MockTodoList, id int) {
				s.EXPECT().GetById(1, id).Return(todo.TodoList{},
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

			list := mock_service.NewMockTodoList(c)
			if listId, ok := testCase.inputId.(int); ok {
				testCase.mockBehavior(list, listId)
			}

			services := &service.Service{TodoList: list}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/api/lists/:id", func(c *gin.Context) {
				c.Set(userCtx, 1)
			}, handler.getListById)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/lists/%v",
				testCase.inputId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.Equal(t, testCase.expectedResponse, w.Body.String())
		})
	}
}

func TestHandler_updateList(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoList, id int,
		input todo.UpdateListInput)

	titleStirng := "Test"
	descriptionString := "Description"
	testTable := []struct {
		name             string
		inputId          any
		inputBody        string
		inputUpdate      todo.UpdateListInput
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:      "OK",
			inputId:   1,
			inputBody: `{"title":"Test","description":"Description"}`,
			inputUpdate: todo.UpdateListInput{
				Title:       &titleStirng,
				Description: &descriptionString,
			},
			mockBehavior: func(s *mock_service.MockTodoList, id int,
				input todo.UpdateListInput) {
				s.EXPECT().Update(1, id, input).Return(nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"status":"ok"}`,
		},
		{
			name:      "Invalid id",
			inputId:   "asd",
			inputBody: `{"title":"Test","description":"Description"}`,
			inputUpdate: todo.UpdateListInput{
				Title:       &titleStirng,
				Description: &descriptionString,
			},
			mockBehavior: func(s *mock_service.MockTodoList, id int,
				input todo.UpdateListInput) {
			},
			expectedStatus:   400,
			expectedResponse: `{"message":"Invalid list id"}`,
		},
		{
			name:      "No list with such id",
			inputId:   10,
			inputBody: `{"title":"Test","description":"Description"}`,
			inputUpdate: todo.UpdateListInput{
				Title:       &titleStirng,
				Description: &descriptionString,
			},
			mockBehavior: func(s *mock_service.MockTodoList, id int,
				input todo.UpdateListInput) {
				s.EXPECT().Update(1, id, input).Return(&todo.ErrNoSuchList{})
			},
			expectedStatus:   200,
			expectedResponse: `{"message":"No list with such id"}`,
		},
		{
			name:      "Invalid request body",
			inputId:   1,
			inputBody: `{}`,
			inputUpdate: todo.UpdateListInput{
				Title:       nil,
				Description: nil,
			},
			mockBehavior: func(s *mock_service.MockTodoList, id int,
				input todo.UpdateListInput) {
				s.EXPECT().Update(1, id, input).Return(
					&todo.ErrInvalidUpdateListInput{})
			},
			expectedStatus:   400,
			expectedResponse: `{"message":"Invalid list update input"}`,
		},
		{
			name:      "Service failure",
			inputId:   1,
			inputBody: `{"title":"Test","description":"Description"}`,
			inputUpdate: todo.UpdateListInput{
				Title:       &titleStirng,
				Description: &descriptionString,
			},
			mockBehavior: func(s *mock_service.MockTodoList, id int,
				input todo.UpdateListInput) {
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

			list := mock_service.NewMockTodoList(c)
			if listId, ok := testCase.inputId.(int); ok {
				testCase.mockBehavior(list, listId, testCase.inputUpdate)
			}

			services := &service.Service{TodoList: list}
			handler := NewHandler(services)

			r := gin.New()
			r.PUT("/api/lists/:id", func(c *gin.Context) {
				c.Set(userCtx, 1)
			}, handler.updateList)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", fmt.Sprintf("/api/lists/%v",
				testCase.inputId), bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.Equal(t, testCase.expectedResponse, w.Body.String())
		})
	}
}

func TestHandler_deleteList(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoList, id int)

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
			mockBehavior: func(s *mock_service.MockTodoList, id int) {
				s.EXPECT().Delete(1, id).Return(nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"status":"ok"}`,
		},
		{
			name:             "Invalid id",
			inputId:          "asd",
			mockBehavior:     func(s *mock_service.MockTodoList, id int) {},
			expectedStatus:   400,
			expectedResponse: `{"message":"Invalid list id"}`,
		},
		{
			name:    "No list with such id",
			inputId: 10,
			mockBehavior: func(s *mock_service.MockTodoList, id int) {
				s.EXPECT().Delete(1, id).Return(nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"status":"ok"}`,
		},
		{
			name:    "Service failure",
			inputId: 1,
			mockBehavior: func(s *mock_service.MockTodoList, id int) {
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

			list := mock_service.NewMockTodoList(c)
			if listId, ok := testCase.inputId.(int); ok {
				testCase.mockBehavior(list, listId)
			}

			services := &service.Service{TodoList: list}
			handler := NewHandler(services)

			r := gin.New()
			r.DELETE("/api/lists/:id", func(c *gin.Context) {
				c.Set(userCtx, 1)
			}, handler.deleteList)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/lists/%v",
				testCase.inputId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.Equal(t, testCase.expectedResponse, w.Body.String())
		})
	}
}

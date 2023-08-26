package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	todo "github.com/OrIX219/todo/pkg"
	"github.com/OrIX219/todo/pkg/service"
	mock_service "github.com/OrIX219/todo/pkg/service/mock"
	_ "github.com/OrIX219/todo/testing"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user todo.User)

	testTable := []struct {
		name             string
		inputBody        string
		inputUser        todo.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:      "OK",
			inputBody: `{"name":"Test","username":"test","password":"qwerty"}`,
			inputUser: todo.User{
				Name:     "Test",
				Username: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"id":1}`,
		},
		{
			name:             "Empty Fields",
			inputBody:        `{"name":"Test","password":"qwerty"}`,
			mockBehavior:     func(s *mock_service.MockAuthorization, user todo.User) {},
			expectedStatus:   400,
			expectedResponse: `{"message":"Invalid request body"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"name":"Test","username":"test","password":"qwerty"}`,
			inputUser: todo.User{
				Name:     "Test",
				Username: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().CreateUser(user).Return(1, errors.New("Service failure"))
			},
			expectedStatus:   500,
			expectedResponse: `{"message":"Service failure"}`,
		}}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.Equal(t, testCase.expectedResponse, w.Body.String())
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, input signInInput)

	testTable := []struct {
		name             string
		inputBody        string
		signInInput      signInInput
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:      "OK",
			inputBody: `{"username":"test","password":"qwerty"}`,
			signInInput: signInInput{
				Username: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, input signInInput) {
				s.EXPECT().GenerateToken(input.Username, input.Password).Return(
					"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTI5MjY1MDEsImlhdCI6MTY5Mjg4MzMwMSwidXNlcl9pZCI6Mn0.5KU4kKRh9Hx-b-hWW46c7XbTJYg3w9l92afZVMgclaQ",
					nil,
				)
			},
			expectedStatus:   200,
			expectedResponse: `{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTI5MjY1MDEsImlhdCI6MTY5Mjg4MzMwMSwidXNlcl9pZCI6Mn0.5KU4kKRh9Hx-b-hWW46c7XbTJYg3w9l92afZVMgclaQ"}`,
		},
		{
			name:             "Empty Fields",
			inputBody:        `{"username":"test"}`,
			mockBehavior:     func(s *mock_service.MockAuthorization, input signInInput) {},
			expectedStatus:   400,
			expectedResponse: `{"message":"Invalid request body"}`,
		},
		{
			name:      "Invalid credentials",
			inputBody: `{"username":"test","password":"ytrewq"}`,
			signInInput: signInInput{
				Username: "test",
				Password: "ytrewq",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, input signInInput) {
				s.EXPECT().GenerateToken(input.Username, input.Password).Return("",
					&todo.ErrInvalidCredentials{})
			},
			expectedStatus:   200,
			expectedResponse: `{"message":"Invalid credentials"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"username":"test","password":"qwerty"}`,
			signInInput: signInInput{
				Username: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, input signInInput) {
				s.EXPECT().GenerateToken(input.Username, input.Password).Return("",
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

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.signInInput)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/sign-in", handler.signIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.Equal(t, testCase.expectedResponse, w.Body.String())
		})
	}
}

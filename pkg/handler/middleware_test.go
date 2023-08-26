package handler

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/OrIX219/todo/pkg/service"
	mock_service "github.com/OrIX219/todo/pkg/service/mock"
	_ "github.com/OrIX219/todo/testing"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestHandler_userIdentity(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, token string)

	testTable := []struct {
		name             string
		headerName       string
		headerValue      string
		token            string
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, nil)
			},
			expectedStatus:   200,
			expectedResponse: "1",
		},
		{
			name:             "No header",
			headerName:       "",
			mockBehavior:     func(s *mock_service.MockAuthorization, token string) {},
			expectedStatus:   401,
			expectedResponse: `{"message":"Empty auth header"}`,
		},
		{
			name:             "Invalid Bearer",
			headerName:       "Authorization",
			headerValue:      "Bearr token",
			mockBehavior:     func(s *mock_service.MockAuthorization, token string) {},
			expectedStatus:   401,
			expectedResponse: `{"message":"Invalid auth header"}`,
		},
		{
			name:             "Invalid token",
			headerName:       "Authorization",
			headerValue:      "Bearer ",
			mockBehavior:     func(s *mock_service.MockAuthorization, token string) {},
			expectedStatus:   401,
			expectedResponse: `{"message":"Empty token"}`,
		},
		{
			name:        "Service failure",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(0, errors.New("Invalid token"))
			},
			expectedStatus:   401,
			expectedResponse: `{"message":"Invalid token"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.token)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/protected", handler.userIdentity, func(c *gin.Context) {
				id, _ := c.Get(userCtx)
				c.String(200, "%d", id.(int))
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.Equal(t, testCase.expectedResponse, w.Body.String())
		})
	}
}

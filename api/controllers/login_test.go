package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kavehjamshidi/fidibo-challenge/domain"
	"github.com/kavehjamshidi/fidibo-challenge/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svcMock := &mocks.LoginService{}
		loginController := NewLoginController(svcMock)

		requestData := domain.LoginRequest{
			Username: "test",
			Password: "test",
		}
		jsonData, err := json.Marshal(requestData)
		assert.NoError(t, err)

		expectedResponse := domain.LoginResponse{
			AccessToken:  "access token",
			RefreshToken: "refresh token",
		}
		expectedJSONResponse, err := json.Marshal(expectedResponse)
		assert.NoError(t, err)

		svcMock.On("Login", requestData).Return(expectedResponse, nil)

		w := httptest.NewRecorder()

		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{Header: make(http.Header)}
		c.Request.Method = http.MethodPost
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonData))

		loginController.Login(c)

		res, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, string(expectedJSONResponse), string(res))
		svcMock.AssertExpectations(t)
	})

	t.Run("invalid body", func(t *testing.T) {
		svcMock := &mocks.LoginService{}
		loginController := NewLoginController(svcMock)

		requestData := domain.LoginRequest{
			Username: "test",
			Password: "",
		}
		jsonData, err := json.Marshal(requestData)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{Header: make(http.Header)}
		c.Request.Method = http.MethodPost
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonData))

		loginController.Login(c)

		res, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		response := domain.ErrorResponse{}
		err = json.Unmarshal(res, &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.NotEmpty(t, response.Message)
		svcMock.AssertExpectations(t)
	})
}

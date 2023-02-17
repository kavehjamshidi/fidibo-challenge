package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kavehjamshidi/fidibo-challenge/domain"
	"github.com/kavehjamshidi/fidibo-challenge/internal/token"
	"github.com/kavehjamshidi/fidibo-challenge/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRefreshToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		secret := "test secret"
		username := "test"
		svcMock := &mocks.RefreshTokenService{}
		refreshTokenController := NewRefreshTokenController(svcMock, secret)

		jwt, err := token.GenerateJWT(username, secret, time.Hour)
		assert.NoError(t, err)

		refreshTokenRequest := domain.RefreshTokenRequest{
			RefreshToken: jwt,
		}
		jsonData, err := json.Marshal(refreshTokenRequest)
		assert.NoError(t, err)

		expectedResponse := domain.RefreshTokenResponse{
			AccessToken:  "access token",
			RefreshToken: "refresh token",
		}
		expectedJSONResponse, err := json.Marshal(expectedResponse)
		assert.NoError(t, err)

		svcMock.On("RefreshToken", username).Return(expectedResponse, nil)

		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{Header: make(http.Header)}
		c.Request.Method = http.MethodPost
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonData))

		refreshTokenController.RefreshToken(c)

		res, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, string(expectedJSONResponse), string(res))
		svcMock.AssertExpectations(t)
	})

	t.Run("invalid token", func(t *testing.T) {
		secret := "test secret"
		svcMock := &mocks.RefreshTokenService{}
		refreshTokenController := NewRefreshTokenController(svcMock, secret)

		jwt := "invalid jwt"

		refreshTokenRequest := domain.RefreshTokenRequest{
			RefreshToken: jwt,
		}
		jsonData, err := json.Marshal(refreshTokenRequest)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{Header: make(http.Header)}
		c.Request.Method = http.MethodPost
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonData))

		refreshTokenController.RefreshToken(c)

		res, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		response := domain.ErrorResponse{}
		err = json.Unmarshal(res, &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.NotEmpty(t, response.Message)
		svcMock.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		secret := "test secret"
		svcMock := &mocks.RefreshTokenService{}
		refreshTokenController := NewRefreshTokenController(svcMock, secret)

		r := gin.Default()
		r.POST("/refresh-token", refreshTokenController.RefreshToken)

		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{Header: make(http.Header)}
		c.Request.Method = http.MethodPost
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = io.NopCloser(bytes.NewBuffer([]byte("invalid body")))

		refreshTokenController.RefreshToken(c)

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

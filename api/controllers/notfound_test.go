package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kavehjamshidi/fidibo-challenge/domain"
	"github.com/stretchr/testify/assert"
)

func TestNotFound(t *testing.T) {
	notFoundController := NewNotFoundController()

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{Header: make(http.Header)}
	c.Request.Method = http.MethodPost

	notFoundController.NotFound(c)

	res, err := io.ReadAll(w.Body)
	assert.NoError(t, err)

	response := domain.ErrorResponse{}
	err = json.Unmarshal(res, &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "Not Found", response.Message)
}

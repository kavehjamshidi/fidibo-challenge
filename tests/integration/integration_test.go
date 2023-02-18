package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kavehjamshidi/fidibo-challenge/api/controllers"
	"github.com/kavehjamshidi/fidibo-challenge/api/routes"
	"github.com/kavehjamshidi/fidibo-challenge/bootstrap"
	"github.com/kavehjamshidi/fidibo-challenge/cache"
	"github.com/kavehjamshidi/fidibo-challenge/db"
	"github.com/kavehjamshidi/fidibo-challenge/domain"
	"github.com/kavehjamshidi/fidibo-challenge/internal/token"
	"github.com/kavehjamshidi/fidibo-challenge/pkg/fidibosearch"
	"github.com/kavehjamshidi/fidibo-challenge/service"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

const (
	fidiboQueryKey  = "q"
	fidiboSearchURL = "https://search.fidibo.com"
)

var (
	router      *gin.Engine
	env         *bootstrap.Env
	redisClient *redis.Client
)

func TestMain(m *testing.M) {
	env = bootstrap.NewEnv()

	redisClient = db.NewRedisClient(context.Background(), env.TestRedisAddress)
	cache := cache.NewCacher(redisClient)

	fidiboClient := fidibosearch.NewFidiboSearcher(fidiboQueryKey, fidiboSearchURL)

	loginSVC := service.NewLoginService(env.AccessTokenExpiry,
		env.AccessTokenSecret,
		env.RefreshTokenExpiry,
		env.RefreshTokenSecret)
	refreshTokenSVC := service.NewRefreshTokenService(env.AccessTokenExpiry,
		env.AccessTokenSecret,
		env.RefreshTokenExpiry,
		env.RefreshTokenSecret)
	searchSVC := service.NewSearchService(cache, fidiboClient)

	loginController := controllers.NewLoginController(loginSVC)
	refreshTokenController := controllers.NewRefreshTokenController(refreshTokenSVC, env.RefreshTokenSecret)
	searchController := controllers.NewSearchController(searchSVC)
	notFoundController := controllers.NewNotFoundController()

	router = gin.Default()

	routes.Setup(router, routes.Controllers{
		SearchController:       searchController,
		LoginController:        loginController,
		RefreshTokenController: refreshTokenController,
	}, env.AccessTokenSecret)

	router.NoRoute(notFoundController.NotFound)

	os.Exit(m.Run())
}

func TestLogin(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		request := domain.LoginRequest{
			Username: "test",
			Password: "test",
		}
		jsonRequest, err := json.Marshal(request)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonRequest))
		assert.NoError(t, err)
		router.ServeHTTP(w, req)

		res, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		response := domain.LoginResponse{}
		err = json.Unmarshal(res, &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, response.AccessToken)
		assert.NotEmpty(t, response.RefreshToken)
	})
}

func TestRefreshToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		jwt, err := token.GenerateJWT("test", env.RefreshTokenSecret, env.RefreshTokenExpiry)
		assert.NoError(t, err)

		request := domain.RefreshTokenRequest{
			RefreshToken: jwt,
		}
		jsonRequest, err := json.Marshal(request)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/refresh-token", bytes.NewReader(jsonRequest))
		assert.NoError(t, err)
		router.ServeHTTP(w, req)

		res, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		response := domain.RefreshTokenResponse{}
		err = json.Unmarshal(res, &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, response.AccessToken)
		assert.NotEmpty(t, response.RefreshToken)
	})

	t.Run("invalid token", func(t *testing.T) {
		request := domain.RefreshTokenRequest{
			RefreshToken: "invalid jwt",
		}
		jsonRequest, err := json.Marshal(request)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/refresh-token", bytes.NewReader(jsonRequest))
		assert.NoError(t, err)
		router.ServeHTTP(w, req)

		res, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		response := domain.ErrorResponse{}
		err = json.Unmarshal(res, &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.NotEmpty(t, response.Message)
	})

	t.Run("expired token", func(t *testing.T) {
		jwt, err := token.GenerateJWT("test", env.RefreshTokenSecret, -time.Minute)
		assert.NoError(t, err)

		request := domain.RefreshTokenRequest{
			RefreshToken: jwt,
		}
		jsonRequest, err := json.Marshal(request)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/refresh-token", bytes.NewReader(jsonRequest))
		assert.NoError(t, err)
		router.ServeHTTP(w, req)

		res, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		response := domain.ErrorResponse{}
		err = json.Unmarshal(res, &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.NotEmpty(t, response.Message)
		assert.Contains(t, response.Message, "token is expired")
	})
}

func TestSearch(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		defer redisClient.FlushAll(context.TODO())

		query := "کافکا"
		jwt, err := token.GenerateJWT("test", env.AccessTokenSecret, env.AccessTokenExpiry)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/search/book?keyword="+query, nil)
		req.Header = make(http.Header)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt))

		assert.NoError(t, err)
		router.ServeHTTP(w, req)

		res, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		response := domain.SearchResult{}
		err = json.Unmarshal(res, &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, response.Books)
	})

	t.Run("success with cache hit", func(t *testing.T) {
		defer redisClient.FlushAll(context.TODO())

		query := "کافکا"
		jwt, err := token.GenerateJWT("test", env.AccessTokenSecret, env.AccessTokenExpiry)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/search/book?keyword="+query, nil)
		req.Header = make(http.Header)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt))

		assert.NoError(t, err)
		router.ServeHTTP(w, req)

		res, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		response := domain.SearchResult{}
		err = json.Unmarshal(res, &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, response.Books)

		w = httptest.NewRecorder()
		req, err = http.NewRequest(http.MethodPost, "/search/book?keyword="+query, nil)
		req.Header = make(http.Header)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt))

		assert.NoError(t, err)
		router.ServeHTTP(w, req)

		res, err = io.ReadAll(w.Body)
		assert.NoError(t, err)

		response = domain.SearchResult{}
		err = json.Unmarshal(res, &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, response.Books)
	})

	t.Run("success with empty query", func(t *testing.T) {
		defer redisClient.FlushAll(context.TODO())

		query := ""
		jwt, err := token.GenerateJWT("test", env.AccessTokenSecret, env.AccessTokenExpiry)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/search/book?keyword="+query, nil)
		req.Header = make(http.Header)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt))

		assert.NoError(t, err)
		router.ServeHTTP(w, req)

		res, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		response := domain.SearchResult{}
		err = json.Unmarshal(res, &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, response.Books)
	})

	t.Run("success with no query", func(t *testing.T) {
		defer redisClient.FlushAll(context.TODO())

		jwt, err := token.GenerateJWT("test", env.AccessTokenSecret, env.AccessTokenExpiry)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/search/book", nil)
		req.Header = make(http.Header)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt))

		assert.NoError(t, err)
		router.ServeHTTP(w, req)

		res, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		response := domain.SearchResult{}
		err = json.Unmarshal(res, &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, response.Books)
	})

	t.Run("test with no auth header", func(t *testing.T) {
		query := "کافکا"

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/search/book?keyword="+query, nil)

		assert.NoError(t, err)
		router.ServeHTTP(w, req)

		res, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		response := domain.ErrorResponse{}
		err = json.Unmarshal(res, &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.NotEmpty(t, response.Message)
	})

	t.Run("test with invalid JWT", func(t *testing.T) {
		query := "کافکا"

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/search/book?keyword="+query, nil)
		req.Header = make(http.Header)
		req.Header.Add("Authorization", "Bearer invalid-token")

		assert.NoError(t, err)
		router.ServeHTTP(w, req)

		res, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		response := domain.ErrorResponse{}
		err = json.Unmarshal(res, &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.NotEmpty(t, response.Message)
	})

	t.Run("test not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/search/books", nil)
		assert.NoError(t, err)
		router.ServeHTTP(w, req)

		res, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		response := domain.ErrorResponse{}
		err = json.Unmarshal(res, &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.NotEmpty(t, response.Message)
		assert.Equal(t, "Not Found", response.Message)
	})
}

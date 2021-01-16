package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jacobsngoodwin/memrizr/account/model"
	"github.com/jacobsngoodwin/memrizr/account/model/apperrors"
	"github.com/jacobsngoodwin/memrizr/account/model/mocks"
	"github.com/stretchr/testify/mock"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func TestSignout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		ctxUser := &model.User{
			UID:   uid,
			Email: "bob1@bob.com",
		}

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// creates a test context for setting a user
		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("user", ctxUser)
		})

		mockTokenService := new(mocks.MockTokenService)
		mockTokenService.On("Signout", mock.AnythingOfType("*context.emptyCtx"), ctxUser.UID).Return(nil)

		NewHandler(&Config{
			R:            router,
			TokenService: mockTokenService,
		})

		request, _ := http.NewRequest(http.MethodPost, "/signout", nil)
		router.ServeHTTP(rr, request)

		respBody, _ := json.Marshal(gin.H{
			"message": "user signed out successfully!",
		})

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
	})

	t.Run("Signout Error", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		ctxUser := &model.User{
			UID:   uid,
			Email: "bob2@bob.com",
		}

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// creates a test context for setting a user
		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("user", ctxUser)
		})

		mockTokenService := new(mocks.MockTokenService)
		mockTokenService.
			On("Signout", mock.AnythingOfType("*context.emptyCtx"), ctxUser.UID).
			Return(apperrors.NewInternal())

		NewHandler(&Config{
			R:            router,
			TokenService: mockTokenService,
		})

		request, _ := http.NewRequest(http.MethodPost, "/signout", nil)
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

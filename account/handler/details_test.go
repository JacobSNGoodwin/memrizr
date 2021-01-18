package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jacobsngoodwin/memrizr/account/model"
	"github.com/jacobsngoodwin/memrizr/account/model/apperrors"
	"github.com/jacobsngoodwin/memrizr/account/model/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDetails(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	uid, _ := uuid.NewRandom()
	ctxUser := &model.User{
		UID: uid,
	}

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("user", ctxUser)
	})

	mockUserService := new(mocks.MockUserService)

	NewHandler(&Config{
		R:           router,
		UserService: mockUserService,
	})

	t.Run("Data binding error", func(t *testing.T) {
		rr := httptest.NewRecorder()

		reqBody, _ := json.Marshal(gin.H{
			"email": "notanemail",
		})
		request, _ := http.NewRequest(http.MethodPut, "/details", bytes.NewBuffer(reqBody))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockUserService.AssertNotCalled(t, "UpdateDetails")
	})

	t.Run("Update success", func(t *testing.T) {
		rr := httptest.NewRecorder()

		newName := "Jacob"
		newEmail := "jacob@jacob.com"
		newWebsite := "https://jacobgoodwin.me"

		reqBody, _ := json.Marshal(gin.H{
			"name":    newName,
			"email":   newEmail,
			"website": newWebsite,
		})

		request, _ := http.NewRequest(http.MethodPut, "/details", bytes.NewBuffer(reqBody))
		request.Header.Set("Content-Type", "application/json")

		userToUpdate := &model.User{
			UID:     ctxUser.UID,
			Name:    newName,
			Email:   newEmail,
			Website: newWebsite,
		}

		updateArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			userToUpdate,
		}

		dbImageURL := "https://jacobgoodwin.me/static/696292a38f493a4283d1a308e4a11732/84d81/Profile.jpg"

		mockUserService.
			On("UpdateDetails", updateArgs...).
			Run(func(args mock.Arguments) {
				userArg := args.Get(1).(*model.User) // arg 0 is context, arg 1 is *User
				userArg.ImageURL = dbImageURL
			}).
			Return(nil)

		router.ServeHTTP(rr, request)

		userToUpdate.ImageURL = dbImageURL
		respBody, _ := json.Marshal(gin.H{
			"user": userToUpdate,
		})

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockUserService.AssertCalled(t, "UpdateDetails", updateArgs...)
	})

	t.Run("Update failure", func(t *testing.T) {
		rr := httptest.NewRecorder()

		newName := "Jacob"
		newEmail := "jacob@jacob.com"
		newWebsite := "https://jacobgoodwin.me"

		reqBody, _ := json.Marshal(gin.H{
			"name":    newName,
			"email":   newEmail,
			"website": newWebsite,
		})

		request, _ := http.NewRequest(http.MethodPut, "/details", bytes.NewBuffer(reqBody))
		request.Header.Set("Content-Type", "application/json")

		userToUpdate := &model.User{
			UID:     ctxUser.UID,
			Name:    newName,
			Email:   newEmail,
			Website: newWebsite,
		}

		updateArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			userToUpdate,
		}

		mockError := apperrors.NewInternal()

		mockUserService.
			On("UpdateDetails", updateArgs...).
			Return(mockError)

		router.ServeHTTP(rr, request)

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockUserService.AssertCalled(t, "UpdateDetails", updateArgs...)
	})
}

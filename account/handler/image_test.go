package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jacobsngoodwin/memrizr/account/model"
	"github.com/jacobsngoodwin/memrizr/account/model/apperrors"
	"github.com/jacobsngoodwin/memrizr/account/model/fixture"
	"github.com/jacobsngoodwin/memrizr/account/model/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestImage(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	uid, _ := uuid.NewRandom()
	ctxUser := model.User{
		UID: uid,
	}

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("user", &ctxUser)
	})

	mockUserService := new(mocks.MockUserService)

	NewHandler(&Config{
		R:            router,
		UserService:  mockUserService,
		MaxBodyBytes: 4 * 1024 * 1024,
	})

	t.Run("Success", func(t *testing.T) {
		rr := httptest.NewRecorder()

		imageURL := "https://www.imageURL.com/1234"

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()

		setProfileImageArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			ctxUser.UID,
			mock.AnythingOfType("*multipart.FileHeader"),
		}

		updatedUser := ctxUser
		updatedUser.ImageURL = imageURL

		mockUserService.On("SetProfileImage", setProfileImageArgs...).Return(&updatedUser, nil)

		request, _ := http.NewRequest(http.MethodPost, "/image", multipartImageFixture.MultipartBody)
		request.Header.Set("Content-Type", multipartImageFixture.ContentType)

		router.ServeHTTP(rr, request)

		respBody, _ := json.Marshal(gin.H{
			"imageUrl": imageURL,
			"message":  "success",
		})

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockUserService.AssertCalled(t, "SetProfileImage", setProfileImageArgs...)
	})

	t.Run("Disallowed mimetype", func(t *testing.T) {
		rr := httptest.NewRecorder()

		multipartImageFixture := fixture.NewMultipartImage("image.txt", "mage/svg+xml")
		defer multipartImageFixture.Close()

		request, _ := http.NewRequest(http.MethodPost, "/image", multipartImageFixture.MultipartBody)
		request.Header.Set("Content-Type", "multipart/form-data")

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		mockUserService.AssertNotCalled(t, "SetProfileImage")
	})

	t.Run("No image file provided", func(t *testing.T) {
		rr := httptest.NewRecorder()

		request, _ := http.NewRequest(http.MethodPost, "/image", nil)
		request.Header.Set("Content-Type", "multipart/form-data")

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		mockUserService.AssertNotCalled(t, "SetProfileImage")
	})

	t.Run("Error from SetProfileImage", func(t *testing.T) {
		// create unique context user for this test
		uid, _ := uuid.NewRandom()
		ctxUser := model.User{
			UID: uid,
		}

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("user", &ctxUser)
		})

		mockUserService := new(mocks.MockUserService)

		NewHandler(&Config{
			R:            router,
			UserService:  mockUserService,
			MaxBodyBytes: 4 * 1024 * 1024,
		})

		rr := httptest.NewRecorder()

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()

		setProfileImageArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			ctxUser.UID,
			mock.AnythingOfType("*multipart.FileHeader"),
		}

		mockError := apperrors.NewInternal()

		mockUserService.On("SetProfileImage", setProfileImageArgs...).Return(nil, mockError)

		request, _ := http.NewRequest(http.MethodPost, "/image", multipartImageFixture.MultipartBody)
		request.Header.Set("Content-Type", multipartImageFixture.ContentType)

		router.ServeHTTP(rr, request)

		assert.Equal(t, apperrors.Status(mockError), rr.Code)

		mockUserService.AssertCalled(t, "SetProfileImage", setProfileImageArgs...)
	})

	// TODO - how to handle large files? Creating large files is very slow
	// maybe create a byte slice and dupe Go into thinking it's an image...?
}

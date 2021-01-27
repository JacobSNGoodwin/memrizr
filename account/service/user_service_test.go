package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jacobsngoodwin/memrizr/account/model"
	"github.com/jacobsngoodwin/memrizr/account/model/apperrors"
	"github.com/jacobsngoodwin/memrizr/account/model/fixture"
	"github.com/jacobsngoodwin/memrizr/account/model/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGet(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUserResp := &model.User{
			UID:   uid,
			Email: "bob@bob.com",
			Name:  "Bobby Bobson",
		}

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})
		mockUserRepository.On("FindByID", mock.Anything, uid).Return(mockUserResp, nil)

		ctx := context.TODO()
		u, err := us.Get(ctx, uid)

		assert.NoError(t, err)
		assert.Equal(t, u, mockUserResp)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})

		mockUserRepository.On("FindByID", mock.Anything, uid).Return(nil, fmt.Errorf("Some error down the call chain"))

		ctx := context.TODO()
		u, err := us.Get(ctx, uid)

		assert.Nil(t, u)
		assert.Error(t, err)
		mockUserRepository.AssertExpectations(t)
	})
}

func TestSignup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUser := &model.User{
			Email:    "bob@bob.com",
			Password: "howdyhoneighbor!",
		}

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})

		// We can use Run method to modify the user when the Create method is called.
		//  We can then chain on a Return method to return no error
		mockUserRepository.
			On("Create", mock.AnythingOfType("*context.emptyCtx"), mockUser).
			Run(func(args mock.Arguments) {
				userArg := args.Get(1).(*model.User) // arg 0 is context, arg 1 is *User
				userArg.UID = uid
			}).Return(nil)

		ctx := context.TODO()
		err := us.Signup(ctx, mockUser)

		assert.NoError(t, err)

		// assert user now has a userID
		assert.Equal(t, uid, mockUser.UID)

		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockUser := &model.User{
			Email:    "bob@bob.com",
			Password: "howdyhoneighbor!",
		}

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})

		mockErr := apperrors.NewConflict("email", mockUser.Email)

		// We can use Run method to modify the user when the Create method is called.
		//  We can then chain on a Return method to return no error
		mockUserRepository.
			On("Create", mock.AnythingOfType("*context.emptyCtx"), mockUser).
			Return(mockErr)

		ctx := context.TODO()
		err := us.Signup(ctx, mockUser)

		// assert error is error we response with in mock
		assert.EqualError(t, err, mockErr.Error())

		mockUserRepository.AssertExpectations(t)
	})
}

func TestSignin(t *testing.T) {
	// setup valid email/pw combo with hashed password to test method
	// response when provided password is invalid
	email := "bob@bob.com"
	validPW := "howdyhoneighbor!"
	hashedValidPW, _ := hashPassword(validPW)
	invalidPW := "howdyhodufus!"

	mockUserRepository := new(mocks.MockUserRepository)
	us := NewUserService(&USConfig{
		UserRepository: mockUserRepository,
	})

	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUser := &model.User{
			Email:    email,
			Password: validPW,
		}

		mockUserResp := &model.User{
			UID:      uid,
			Email:    email,
			Password: hashedValidPW,
		}

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			email,
		}

		// We can use Run method to modify the user when the Create method is called.
		//  We can then chain on a Return method to return no error
		mockUserRepository.
			On("FindByEmail", mockArgs...).Return(mockUserResp, nil)

		ctx := context.TODO()
		err := us.Signin(ctx, mockUser)

		assert.NoError(t, err)
		mockUserRepository.AssertCalled(t, "FindByEmail", mockArgs...)
	})

	t.Run("Invalid email/password combination", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUser := &model.User{
			Email:    email,
			Password: invalidPW,
		}

		mockUserResp := &model.User{
			UID:      uid,
			Email:    email,
			Password: hashedValidPW,
		}

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			email,
		}

		// We can use Run method to modify the user when the Create method is called.
		//  We can then chain on a Return method to return no error
		mockUserRepository.
			On("FindByEmail", mockArgs...).Return(mockUserResp, nil)

		ctx := context.TODO()
		err := us.Signin(ctx, mockUser)

		assert.Error(t, err)
		assert.EqualError(t, err, "Invalid email and password combination")
		mockUserRepository.AssertCalled(t, "FindByEmail", mockArgs...)
	})
}

func TestUpdateDetails(t *testing.T) {
	mockUserRepository := new(mocks.MockUserRepository)
	us := NewUserService(&USConfig{
		UserRepository: mockUserRepository,
	})

	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUser := &model.User{
			UID:     uid,
			Email:   "new@bob.com",
			Website: "https://jacobgoodwin.me",
			Name:    "A New Bob!",
		}

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUser,
		}

		mockUserRepository.
			On("Update", mockArgs...).Return(nil)

		ctx := context.TODO()
		err := us.UpdateDetails(ctx, mockUser)

		assert.NoError(t, err)
		mockUserRepository.AssertCalled(t, "Update", mockArgs...)
	})

	t.Run("Failure", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUser := &model.User{
			UID: uid,
		}

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUser,
		}

		mockError := apperrors.NewInternal()

		mockUserRepository.
			On("Update", mockArgs...).Return(mockError)

		ctx := context.TODO()
		err := us.UpdateDetails(ctx, mockUser)
		assert.Error(t, err)

		apperror, ok := err.(*apperrors.Error)
		assert.True(t, ok)
		assert.Equal(t, apperrors.Internal, apperror.Type)

		mockUserRepository.AssertCalled(t, "Update", mockArgs...)
	})
}

func TestSetProfileImage(t *testing.T) {
	mockUserRepository := new(mocks.MockUserRepository)
	mockImageRepository := new(mocks.MockImageRepository)

	us := NewUserService(&USConfig{
		UserRepository:  mockUserRepository,
		ImageRepository: mockImageRepository,
	})

	t.Run("Successful new image", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		// does not have have imageURL
		mockUser := &model.User{
			UID:     uid,
			Email:   "new@bob.com",
			Website: "https://jacobgoodwin.me",
			Name:    "A New Bob!",
		}

		findByIDArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			uid,
		}
		mockUserRepository.On("FindByID", findByIDArgs...).Return(mockUser, nil)

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()
		imageFileHeader := multipartImageFixture.GetFormFile()
		imageFile, _ := imageFileHeader.Open()

		updateProfileArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("string"),
			imageFile,
		}

		imageURL := "http://imageurl.com/jdfkj34kljl"

		mockImageRepository.
			On("UpdateProfile", updateProfileArgs...).
			Return(imageURL, nil)

		updateImageArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUser.UID,
			imageURL,
		}

		mockUpdatedUser := &model.User{
			UID:      uid,
			Email:    "new@bob.com",
			Website:  "https://jacobgoodwin.me",
			Name:     "A New Bob!",
			ImageURL: imageURL,
		}

		mockUserRepository.
			On("UpdateImage", updateImageArgs...).
			Return(mockUpdatedUser, nil)

		ctx := context.TODO()

		updatedUser, err := us.SetProfileImage(ctx, mockUser.UID, imageFileHeader)

		assert.NoError(t, err)
		assert.Equal(t, mockUpdatedUser, updatedUser)
		mockUserRepository.AssertCalled(t, "FindByID", findByIDArgs...)
		mockImageRepository.AssertCalled(t, "UpdateProfile", updateProfileArgs...)
		mockUserRepository.AssertCalled(t, "UpdateImage", updateImageArgs...)
	})

	t.Run("Successful update image", func(t *testing.T) {
		uid, _ := uuid.NewRandom()
		imageURL := "http://imageurl.com/jdfkj34kljl"

		// has imageURL
		mockUser := &model.User{
			UID:      uid,
			Email:    "new@bob.com",
			Website:  "https://jacobgoodwin.me",
			Name:     "A New Bob!",
			ImageURL: imageURL,
		}

		findByIDArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			uid,
		}
		mockUserRepository.On("FindByID", findByIDArgs...).Return(mockUser, nil)

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()
		imageFileHeader := multipartImageFixture.GetFormFile()
		imageFile, _ := imageFileHeader.Open()

		updateProfileArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("string"),
			imageFile,
		}

		mockImageRepository.
			On("UpdateProfile", updateProfileArgs...).
			Return(imageURL, nil)

		updateImageArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUser.UID,
			imageURL,
		}

		mockUpdatedUser := &model.User{
			UID:      uid,
			Email:    "new@bob.com",
			Website:  "https://jacobgoodwin.me",
			Name:     "A New Bob!",
			ImageURL: imageURL,
		}

		mockUserRepository.
			On("UpdateImage", updateImageArgs...).
			Return(mockUpdatedUser, nil)

		ctx := context.TODO()

		updatedUser, err := us.SetProfileImage(ctx, uid, imageFileHeader)

		assert.NoError(t, err)
		assert.Equal(t, mockUpdatedUser, updatedUser)
		mockUserRepository.AssertCalled(t, "FindByID", findByIDArgs...)
		mockImageRepository.AssertCalled(t, "UpdateProfile", updateProfileArgs...)
		mockUserRepository.AssertCalled(t, "UpdateImage", updateImageArgs...)
	})

	t.Run("UserRepository FindByID Error", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		findByIDArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			uid,
		}
		mockError := apperrors.NewInternal()
		mockUserRepository.On("FindByID", findByIDArgs...).Return(nil, mockError)

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()
		imageFileHeader := multipartImageFixture.GetFormFile()

		ctx := context.TODO()

		updatedUser, err := us.SetProfileImage(ctx, uid, imageFileHeader)

		assert.Error(t, err)
		assert.Nil(t, updatedUser)
		mockUserRepository.AssertCalled(t, "FindByID", findByIDArgs...)
		mockImageRepository.AssertNotCalled(t, "UpdateProfile")
		mockUserRepository.AssertNotCalled(t, "UpdateImage")
	})

	t.Run("ImageRepository Error", func(t *testing.T) {
		// need to create a new UserService and repository
		// because testify has no way to overwrite a mock's
		// "On" call.
		mockUserRepository := new(mocks.MockUserRepository)
		mockImageRepository := new(mocks.MockImageRepository)

		us := NewUserService(&USConfig{
			UserRepository:  mockUserRepository,
			ImageRepository: mockImageRepository,
		})

		uid, _ := uuid.NewRandom()
		imageURL := "http://imageurl.com/jdfkj34kljl"

		// has imageURL
		mockUser := &model.User{
			UID:      uid,
			Email:    "new@bob.com",
			Website:  "https://jacobgoodwin.me",
			Name:     "A New Bob!",
			ImageURL: imageURL,
		}

		findByIDArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			uid,
		}
		mockUserRepository.On("FindByID", findByIDArgs...).Return(mockUser, nil)

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()
		imageFileHeader := multipartImageFixture.GetFormFile()
		imageFile, _ := imageFileHeader.Open()

		updateProfileArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("string"),
			imageFile,
		}

		mockError := apperrors.NewInternal()
		mockImageRepository.
			On("UpdateProfile", updateProfileArgs...).
			Return(nil, mockError)

		ctx := context.TODO()
		updatedUser, err := us.SetProfileImage(ctx, uid, imageFileHeader)

		assert.Nil(t, updatedUser)
		assert.Error(t, err)
		mockUserRepository.AssertCalled(t, "FindByID", findByIDArgs...)
		mockImageRepository.AssertCalled(t, "UpdateProfile", updateProfileArgs...)
		mockUserRepository.AssertNotCalled(t, "UpdateImage")
	})

	t.Run("UserRepository UpdateImage Error", func(t *testing.T) {
		uid, _ := uuid.NewRandom()
		imageURL := "http://imageurl.com/jdfkj34kljl"

		// has imageURL
		mockUser := &model.User{
			UID:      uid,
			Email:    "new@bob.com",
			Website:  "https://jacobgoodwin.me",
			Name:     "A New Bob!",
			ImageURL: imageURL,
		}

		findByIDArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			uid,
		}
		mockUserRepository.On("FindByID", findByIDArgs...).Return(mockUser, nil)

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()
		imageFileHeader := multipartImageFixture.GetFormFile()
		imageFile, _ := imageFileHeader.Open()

		updateProfileArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("string"),
			imageFile,
		}

		mockImageRepository.
			On("UpdateProfile", updateProfileArgs...).
			Return(imageURL, nil)

		updateImageArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUser.UID,
			imageURL,
		}

		mockError := apperrors.NewInternal()
		mockUserRepository.
			On("UpdateImage", updateImageArgs...).
			Return(nil, mockError)

		ctx := context.TODO()

		updatedUser, err := us.SetProfileImage(ctx, uid, imageFileHeader)

		assert.Error(t, err)
		assert.Nil(t, updatedUser)
		mockImageRepository.AssertCalled(t, "UpdateProfile", updateProfileArgs...)
		mockUserRepository.AssertCalled(t, "UpdateImage", updateImageArgs...)
	})

	// TODO - Create non image file for test
}

package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/jacobsngoodwin/memrizr/account/model"
	"github.com/jacobsngoodwin/memrizr/account/model/apperrors"
	"github.com/jacobsngoodwin/memrizr/account/model/mocks"
)

func TestNewPairFromUser(t *testing.T) {
	var idExp int64 = 15 * 60
	var refreshExp int64 = 3 * 24 * 2600
	priv, _ := ioutil.ReadFile("../rsa_private_test.pem")
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)
	if err != nil {
		privKey, _ = generatePrivateKey(2048)
	}
	pub, _ := ioutil.ReadFile("../rsa_public_test.pem")
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		pubKey = &privKey.PublicKey
	}
	secret := "anotsorandomtestsecret"

	mockTokenRepository := new(mocks.MockTokenRepository)

	// instantiate a common token service to be used by all tests
	tokenService := NewTokenService(&TSConfig{
		TokenRepository:       mockTokenRepository,
		PrivKey:               privKey,
		PubKey:                pubKey,
		RefreshSecret:         secret,
		IDExpirationSecs:      idExp,
		RefreshExpirationSecs: refreshExp,
	})

	// include password to make sure it is not serialized
	// since json tag is "-"
	uid, _ := uuid.NewRandom()
	u := &model.User{
		UID:      uid,
		Email:    "bob@bob.com",
		Password: "blarghedymcblarghface",
	}

	// Setup mock call responses in setup before t.Run statements
	uidErrorCase, _ := uuid.NewRandom()
	uErrorCase := &model.User{
		UID:      uidErrorCase,
		Email:    "failure@failure.com",
		Password: "blarghedymcblarghface",
	}
	prevID := "a_previous_tokenID"

	setSuccessArguments := mock.Arguments{
		mock.AnythingOfType("*context.emptyCtx"),
		u.UID.String(),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Duration"),
	}

	setErrorArguments := mock.Arguments{
		mock.AnythingOfType("*context.emptyCtx"),
		uidErrorCase.String(),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Duration"),
	}

	deleteWithPrevIDArguments := mock.Arguments{
		mock.AnythingOfType("*context.emptyCtx"),
		u.UID.String(),
		prevID,
	}

	// mock call argument/responses
	mockTokenRepository.On("SetRefreshToken", setSuccessArguments...).Return(nil)
	mockTokenRepository.On("SetRefreshToken", setErrorArguments...).Return(fmt.Errorf("Error setting refresh token"))
	mockTokenRepository.On("DeleteRefreshToken", deleteWithPrevIDArguments...).Return(nil)

	t.Run("Returns a token pair with proper values", func(t *testing.T) {
		ctx := context.Background()                                    // updated from context.TODO()
		tokenPair, err := tokenService.NewPairFromUser(ctx, u, prevID) // replaced "" with prevID from setup
		assert.NoError(t, err)

		// SetRefreshToken should be called with setSuccessArguments
		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setSuccessArguments...)
		// DeleteRefreshToken should not be called since prevID is ""
		mockTokenRepository.AssertCalled(t, "DeleteRefreshToken", deleteWithPrevIDArguments...)

		var s string
		assert.IsType(t, s, tokenPair.IDToken.SS)

		// decode the Base64URL encoded string
		// simpler to use jwt library which is already imported
		idTokenClaims := &idTokenCustomClaims{}

		_, err = jwt.ParseWithClaims(tokenPair.IDToken.SS, idTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return pubKey, nil
		})

		assert.NoError(t, err)

		// assert claims on idToken
		expectedClaims := []interface{}{
			u.UID,
			u.Email,
			u.Name,
			u.ImageURL,
			u.Website,
		}
		actualIDClaims := []interface{}{
			idTokenClaims.User.UID,
			idTokenClaims.User.Email,
			idTokenClaims.User.Name,
			idTokenClaims.User.ImageURL,
			idTokenClaims.User.Website,
		}

		assert.ElementsMatch(t, expectedClaims, actualIDClaims)
		assert.Empty(t, idTokenClaims.User.Password) // password should never be encoded to json

		expiresAt := time.Unix(idTokenClaims.StandardClaims.ExpiresAt, 0)
		expectedExpiresAt := time.Now().Add(time.Duration(idExp) * time.Second)
		assert.WithinDuration(t, expectedExpiresAt, expiresAt, 5*time.Second)

		refreshTokenClaims := &refreshTokenCustomClaims{}
		_, err = jwt.ParseWithClaims(tokenPair.RefreshToken.SS, refreshTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		assert.IsType(t, s, tokenPair.RefreshToken.SS)

		// assert claims on refresh token
		assert.NoError(t, err)
		assert.Equal(t, u.UID, refreshTokenClaims.UID)

		expiresAt = time.Unix(refreshTokenClaims.StandardClaims.ExpiresAt, 0)
		expectedExpiresAt = time.Now().Add(time.Duration(refreshExp) * time.Second)
		assert.WithinDuration(t, expectedExpiresAt, expiresAt, 5*time.Second)
	})

	t.Run("Error setting refresh token", func(t *testing.T) {
		ctx := context.Background()
		_, err := tokenService.NewPairFromUser(ctx, uErrorCase, "")
		assert.Error(t, err) // should return an error

		// SetRefreshToken should be called with setErrorArguments
		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setErrorArguments...)
		// DeleteRefreshToken should not be since SetRefreshToken causes method to return
		mockTokenRepository.AssertNotCalled(t, "DeleteRefreshToken")
	})

	t.Run("Empty string provided for prevID", func(t *testing.T) {
		ctx := context.Background()
		_, err := tokenService.NewPairFromUser(ctx, u, "")
		assert.NoError(t, err)

		// SetRefreshToken should be called with setSuccessArguments
		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setSuccessArguments...)
		// DeleteRefreshToken should not be called since prevID is ""
		mockTokenRepository.AssertNotCalled(t, "DeleteRefreshToken")
	})

	t.Run("Prev token not in repository", func(t *testing.T) {
		ctx := context.Background()
		uid, _ := uuid.NewRandom()
		u := &model.User{
			UID: uid,
		}

		tokenIDNotInRepo := "not_in_token_repo"

		deleteArgs := mock.Arguments{
			ctx,
			u.UID.String(),
			tokenIDNotInRepo,
		}

		mockError := apperrors.NewAuthorization("Invalid refresh token")
		mockTokenRepository.
			On("DeleteRefreshToken", deleteArgs...).
			Return(mockError)

		_, err := tokenService.NewPairFromUser(ctx, u, tokenIDNotInRepo)
		assert.Error(t, err)

		appError, ok := err.(*apperrors.Error)

		assert.True(t, ok)
		assert.Equal(t, apperrors.Authorization, appError.Type)
		mockTokenRepository.AssertCalled(t, "DeleteRefreshToken", deleteArgs...)
		mockTokenRepository.AssertNotCalled(t, "SetRefreshToken")
	})
}

func TestSignout(t *testing.T) {
	mockTokenRepository := new(mocks.MockTokenRepository)
	tokenService := NewTokenService(&TSConfig{
		TokenRepository: mockTokenRepository,
	})

	t.Run("No error", func(t *testing.T) {
		uidSuccess, _ := uuid.NewRandom()
		mockTokenRepository.
			On("DeleteUserRefreshTokens", mock.AnythingOfType("*context.emptyCtx"), uidSuccess.String()).
			Return(nil)

		ctx := context.Background()
		err := tokenService.Signout(ctx, uidSuccess)
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		uidError, _ := uuid.NewRandom()
		mockTokenRepository.
			On("DeleteUserRefreshTokens", mock.AnythingOfType("*context.emptyCtx"), uidError.String()).
			Return(apperrors.NewInternal())

		ctx := context.Background()
		err := tokenService.Signout(ctx, uidError)

		assert.Error(t, err)

		apperr, ok := err.(*apperrors.Error)
		assert.True(t, ok)
		assert.Equal(t, apperr.Type, apperrors.Internal)
	})
}

func TestValidateIDToken(t *testing.T) {
	var idExp int64 = 15 * 60

	priv, _ := ioutil.ReadFile("../rsa_private_test.pem")
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)
	if err != nil {
		privKey, _ = generatePrivateKey(2048)
	}
	pub, _ := ioutil.ReadFile("../rsa_public_test.pem")
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		pubKey = &privKey.PublicKey
	}

	// instantiate a common token service to be used by all tests
	tokenService := NewTokenService(&TSConfig{
		PrivKey:          privKey,
		PubKey:           pubKey,
		IDExpirationSecs: idExp,
	})

	// include password to make sure it is not serialized
	// since json tag is "-"
	uid, _ := uuid.NewRandom()
	u := &model.User{
		UID:      uid,
		Email:    "bob@bob.com",
		Password: "blarghedymcblarghface",
	}

	t.Run("Valid token", func(t *testing.T) {
		// maybe not the best approach to depend on utility method
		// token will be valid for 15 minutes
		ss, _ := generateIDToken(u, privKey, idExp)

		uFromToken, err := tokenService.ValidateIDToken(ss)
		assert.NoError(t, err)

		assert.ElementsMatch(
			t,
			[]interface{}{u.Email, u.Name, u.UID, u.Website, u.ImageURL},
			[]interface{}{uFromToken.Email, uFromToken.Name, uFromToken.UID, uFromToken.Website, uFromToken.ImageURL},
		)
	})

	t.Run("Expired token", func(t *testing.T) {
		// maybe not the best approach to depend on utility method
		// token will be valid for 15 minutes
		ss, _ := generateIDToken(u, privKey, -1) // expires one second ago

		expectedErr := apperrors.NewAuthorization("Unable to verify user from idToken")

		_, err := tokenService.ValidateIDToken(ss)
		assert.EqualError(t, err, expectedErr.Message)
	})

	t.Run("Invalid signature", func(t *testing.T) {
		// maybe not the best approach to depend on utility method
		// token will be valid for 15 minutes
		ss, _ := generateIDToken(u, privKey, -1) // expires one second ago

		expectedErr := apperrors.NewAuthorization("Unable to verify user from idToken")

		_, err := tokenService.ValidateIDToken(ss)
		assert.EqualError(t, err, expectedErr.Message)
	})

	// TODO - Add other invalid token types
}

func TestValidateRefreshToken(t *testing.T) {
	var refreshExp int64 = 3 * 24 * 2600
	secret := "anotsorandomtestsecret"

	tokenService := NewTokenService(&TSConfig{
		RefreshSecret:         secret,
		RefreshExpirationSecs: refreshExp,
	})

	uid, _ := uuid.NewRandom()
	u := &model.User{
		UID:      uid,
		Email:    "bob@bob.com",
		Password: "blarghedymcblarghface",
	}

	t.Run("Valid token", func(t *testing.T) {
		testRefreshToken, _ := generateRefreshToken(u.UID, secret, refreshExp)

		validatedRefreshToken, err := tokenService.ValidateRefreshToken(testRefreshToken.SS)
		assert.NoError(t, err)

		assert.Equal(t, u.UID, validatedRefreshToken.UID)
		assert.Equal(t, testRefreshToken.SS, validatedRefreshToken.SS)
		assert.Equal(t, u.UID, validatedRefreshToken.UID)
	})

	t.Run("Expired token", func(t *testing.T) {
		testRefreshToken, _ := generateRefreshToken(u.UID, secret, -1)

		expectedErr := apperrors.NewAuthorization("Unable to verify user from refresh token")

		_, err := tokenService.ValidateRefreshToken(testRefreshToken.SS)
		assert.EqualError(t, err, expectedErr.Message)
	})
}

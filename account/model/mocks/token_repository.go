package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

// ... imports omitted

// MockTokenRepository is a mock type for model.TokenRepository
type MockTokenRepository struct {
	mock.Mock
}

// SetRefreshToken is a mock of model.TokenRepository SetRefreshToken
func (m *MockTokenRepository) SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error {
	ret := m.Called(ctx, userID, tokenID, expiresIn)

	var r0 error

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

// DeleteRefreshToken is a mock of model.TokenRepository DeleteRefreshToken
func (m *MockTokenRepository) DeleteRefreshToken(ctx context.Context, userID string, prevTokenID string) error {
	ret := m.Called(ctx, userID, prevTokenID)

	var r0 error

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

// DeleteUserRefreshTokens mocks concrete DeleteUserRefreshToken
func (m *MockTokenRepository) DeleteUserRefreshTokens(ctx context.Context, userID string) error {
	ret := m.Called(ctx, userID)

	var r0 error

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

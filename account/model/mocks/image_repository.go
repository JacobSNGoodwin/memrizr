package mocks

import (
	"context"
	"mime/multipart"

	"github.com/stretchr/testify/mock"
)

// MockImageRepository is a mock type for model.ImageRepository
type MockImageRepository struct {
	mock.Mock
}

// UpdateProfile is mock of representations of ImageRepository Update Profile
func (m *MockImageRepository) UpdateProfile(ctx context.Context, objName string, imageFile multipart.File) (string, error) {
	// args that will be passed to "Return" in the tests, when function
	// is called with a uid. Hence the name "ret"
	ret := m.Called(ctx, objName, imageFile)

	// first value passed to "Return"
	var r0 string
	if ret.Get(0) != nil {
		// we can just return this if we know we won't be passing function to "Return"
		r0 = ret.Get(0).(string)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

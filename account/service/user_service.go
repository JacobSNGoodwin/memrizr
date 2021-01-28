package service

import (
	"context"
	"log"
	"mime/multipart"
	"net/url"
	"path"

	"github.com/google/uuid"
	"github.com/jacobsngoodwin/memrizr/account/model"
	"github.com/jacobsngoodwin/memrizr/account/model/apperrors"
)

// userService acts as a struct for injecting an implementation of UserRepository
// for use in service methods
type userService struct {
	UserRepository  model.UserRepository
	ImageRepository model.ImageRepository
}

// USConfig will hold repositories that will eventually be injected into this
// this service layer
type USConfig struct {
	UserRepository  model.UserRepository
	ImageRepository model.ImageRepository
}

// NewUserService is a factory function for
// initializing a UserService with its repository layer dependencies
func NewUserService(c *USConfig) model.UserService {
	return &userService{
		UserRepository:  c.UserRepository,
		ImageRepository: c.ImageRepository,
	}
}

func (s *userService) ClearProfileImage(
	ctx context.Context,
	uid uuid.UUID,
) error {
	user, err := s.UserRepository.FindByID(ctx, uid)

	if err != nil {
		return err
	}

	if user.ImageURL == "" {
		return nil
	}

	objName, err := objNameFromURL(user.ImageURL)
	if err != nil {
		return err
	}

	err = s.ImageRepository.DeleteProfile(ctx, objName)
	if err != nil {
		return err
	}

	_, err = s.UserRepository.UpdateImage(ctx, uid, "")

	if err != nil {
		return err
	}

	return nil
}

// Get retrieves a user based on their uuid
func (s *userService) Get(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	u, err := s.UserRepository.FindByID(ctx, uid)

	return u, err
}

// Signup reaches our to a UserRepository to verify the
// email address is available and signs up the user if this is the case
func (s *userService) Signup(ctx context.Context, u *model.User) error {
	pw, err := hashPassword(u.Password)

	if err != nil {
		log.Printf("Unable to signup user for email: %v\n", u.Email)
		return apperrors.NewInternal()
	}

	// now I realize why I originally used Signup(ctx, email, password)
	// then created a user. It's somewhat un-natural to mutate the user here
	u.Password = pw

	if err := s.UserRepository.Create(ctx, u); err != nil {
		return err
	}

	// If we get around to adding events, we'd Publish it here
	// err := s.EventsBroker.PublishUserUpdated(u, true)

	// if err != nil {
	// 	return nil, apperrors.NewInternal()
	// }

	return nil
}

// Signin reaches our to a UserRepository check if the user exists
// and then compares the supplied password with the provided password
// if a valid email/password combo is provided, u will hold all
// available user fields
func (s *userService) Signin(ctx context.Context, u *model.User) error {
	uFetched, err := s.UserRepository.FindByEmail(ctx, u.Email)

	// Will return NotAuthorized to client to omit details of why
	if err != nil {
		return apperrors.NewAuthorization("Invalid email and password combination")
	}

	// verify password - we previously created this method
	match, err := comparePasswords(uFetched.Password, u.Password)

	if err != nil {
		return apperrors.NewInternal()
	}

	if !match {
		return apperrors.NewAuthorization("Invalid email and password combination")
	}

	*u = *uFetched
	return nil
}

func (s *userService) UpdateDetails(ctx context.Context, u *model.User) error {
	// Update user in UserRepository
	err := s.UserRepository.Update(ctx, u)

	if err != nil {
		return err
	}

	// // Publish user updated
	// err = s.EventsBroker.PublishUserUpdated(u, false)
	// if err != nil {
	// 	return apperrors.NewInternal()
	// }

	return nil
}

func (s *userService) SetProfileImage(
	ctx context.Context,
	uid uuid.UUID,
	imageFileHeader *multipart.FileHeader,
) (*model.User, error) {
	u, err := s.UserRepository.FindByID(ctx, uid)

	if err != nil {
		return nil, err
	}

	objName, err := objNameFromURL(u.ImageURL)

	if err != nil {
		return nil, err
	}

	imageFile, err := imageFileHeader.Open()
	if err != nil {
		log.Printf("Failed to open image file: %v\n", err)
		return nil, apperrors.NewInternal()
	}

	// Upload user's image to ImageRepository
	// Possibly received updated imageURL
	imageURL, err := s.ImageRepository.UpdateProfile(ctx, objName, imageFile)

	if err != nil {
		log.Printf("Unable to upload image to cloud provider: %v\n", err)
		return nil, err
	}

	updatedUser, err := s.UserRepository.UpdateImage(ctx, u.UID, imageURL)

	if err != nil {
		log.Printf("Unable to update imageURL: %v\n", err)
		return nil, err
	}

	return updatedUser, nil
}

func objNameFromURL(imageURL string) (string, error) {
	// if user doesn't have imageURL - create one
	// otherwise, extract last part of URL to get cloud storage object name
	if imageURL == "" {
		objID, _ := uuid.NewRandom()
		return objID.String(), nil
	}

	// split off last part of URL, which is the image's storage object ID
	urlPath, err := url.Parse(imageURL)

	if err != nil {
		log.Printf("Failed to parse objectName from imageURL: %v\n", imageURL)
		return "", apperrors.NewInternal()
	}

	// get "path" of url (everything after domain)
	// then get "base", the last part
	return path.Base(urlPath.Path), nil
}

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/CLoouis/image-uploader/pkg/api/user"
	"github.com/CLoouis/image-uploader/pkg/utl/hash"
)

type (
	UserServiceImpl struct {
		userRepository user.UserRepository
		timeout        time.Duration
	}
)

func NewUserServiceImpl(userRepository user.UserRepository, timeout time.Duration) user.UserService {
	return &UserServiceImpl{
		userRepository: userRepository,
		timeout:        timeout,
	}
}

func (u *UserServiceImpl) Create(ctx context.Context, userData user.User) (string, error) {
	password := userData.Password
	userData.Password = hash.HashPassword(password)

	result, err := u.userRepository.Create(ctx, userData)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (u *UserServiceImpl) Me(ctx context.Context) (user.UserInformation, error) {
	userID := fmt.Sprintf("%v", ctx.Value("id"))

	userData, err := u.userRepository.FindById(ctx, userID)
	if err != nil {
		return user.UserInformation{}, err
	}

	return userData.GetUserInformation(), nil
}

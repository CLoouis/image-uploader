package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CLoouis/image-uploader/pkg/api/image"
	"github.com/CLoouis/image-uploader/pkg/api/user"
	"github.com/CLoouis/image-uploader/pkg/utl/uploader"
	"github.com/google/uuid"
)

type (
	ImageServiceImpl struct {
		imageRepository image.ImageRepository
		userRepository  user.UserRepository
		timeout         time.Duration
		uploader        uploader.Uploader
	}
)

func NewImageServiceImpl(imageRepository image.ImageRepository, userRepository user.UserRepository, timeout time.Duration, uploader uploader.Uploader) image.ImageService {
	return &ImageServiceImpl{
		imageRepository: imageRepository,
		userRepository:  userRepository,
		timeout:         timeout,
		uploader:        uploader,
	}
}

func (svc ImageServiceImpl) SaveImageMetadata(c context.Context, image image.Image) error {
	ctx, cancel := context.WithTimeout(c, svc.timeout)
	defer cancel()

	image.Timestamp = time.Now()
	image.UserId = fmt.Sprintf("%v", c.Value("id"))

	err := svc.imageRepository.SaveImageMetadata(ctx, image)
	if err != nil {
		return err
	}

	return nil
}

func (svc ImageServiceImpl) HandleGetPresignUploadUrlRequest(c context.Context) (image.URLResponse, error) {
	fileName := uuid.New().String()
	url, err := svc.uploader.GetPresignUploadUrl(fileName)
	if err != nil {
		return image.URLResponse{}, err
	}
	return image.URLResponse{FileName: fileName, URL: url}, nil
}

func (svc ImageServiceImpl) GetImageByFileName(c context.Context, fileName string) (image.ImageResponse, error) {
	ctx, cancel := context.WithTimeout(c, svc.timeout)
	defer cancel()

	imageData, err := svc.imageRepository.GetImageMetadataByFileName(ctx, fileName)
	if err != nil {
		return image.ImageResponse{}, err
	}

	userId := fmt.Sprintf("%v", c.Value("id"))

	if imageData.UserId != userId {
		return image.ImageResponse{}, errors.New("unauthorized")
	}

	url, err := svc.uploader.GetPresignFetchUrl(fileName)
	if err != nil {
		return image.ImageResponse{}, err
	}

	result := image.ImageResponse{
		Id:        imageData.Id,
		URL:       url,
		UserId:    imageData.UserId,
		Timestamp: imageData.Timestamp,
	}

	return result, nil
}

func (svc ImageServiceImpl) GetImageByUserId(c context.Context, userId string) (image.ImageResponse, error) {
	panic("")
}

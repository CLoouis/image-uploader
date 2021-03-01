package repository

import (
	"context"
	"errors"

	"github.com/CLoouis/image-uploader/pkg/api/image"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	ImageRepositoryImpl struct {
		imageCollection *mongo.Collection
	}
)

func NewImageRepositoryImpl(imageCollection *mongo.Collection) image.ImageRepository {
	return &ImageRepositoryImpl{imageCollection: imageCollection}
}

func (i ImageRepositoryImpl) SaveImageMetadata(c context.Context, imageData image.Image) error {
	_, err := i.imageCollection.InsertOne(c, imageData)
	if err != nil {
		return err
	}

	return nil
}

func (i ImageRepositoryImpl) GetImageMetadataByFileName(c context.Context, fileName string) (image.Image, error) {
	var imageData image.Image
	err := i.imageCollection.FindOne(c, bson.M{"filename": fileName}).Decode(&imageData)
	if err != nil {
		return image.Image{}, err
	}

	return imageData, nil
}

func (i ImageRepositoryImpl) GetImageMetadataByUserId(c context.Context, userId string) ([]image.Image, error) {
	var listOfImage []image.Image
	cursor, err := i.imageCollection.Find(c, bson.M{"user_id": userId})
	if err != nil {
		return []image.Image{}, err
	}

	for cursor.Next(c) {
		var imageData image.Image
		if err = cursor.Decode(&imageData); err != nil {
			return listOfImage, err
		}

		listOfImage = append(listOfImage, imageData)
	}

	if err := cursor.Err(); err != nil {
		return listOfImage, err
	}

	cursor.Close(c)
	if len(listOfImage) == 0 {
		return listOfImage, errors.New("not found")
	}

	return listOfImage, nil
}

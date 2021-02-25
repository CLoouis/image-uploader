package image

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Image struct {
		Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
		FileName  string             `json:"filename" bson:"filename"`
		UserId    string             `json:"user_id" bson:"user_id"`
		Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
	}

	ImageResponse struct {
		Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
		URL       string             `json:"url" bson:"url"`
		UserId    string             `json:"user_id" bson:"user_id"`
		Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
	}

	URLResponse struct {
		FileName string `json:"filename"`
		URL      string `json:"url"`
	}

	ImageRepository interface {
		SaveImageMetadata(context.Context, Image) error
		GetImageMetadataByFileName(context.Context, string) (Image, error)
		GetImageMetadataByUserId(context.Context, string) (Image, error)
	}

	ImageService interface {
		SaveImageMetadata(context.Context, Image) error
		HandleGetPresignUploadUrlRequest(context.Context) (URLResponse, error)
		GetImageByFileName(context.Context, string) (ImageResponse, error)
		GetImageByUserId(context.Context, string) (ImageResponse, error)
	}
)

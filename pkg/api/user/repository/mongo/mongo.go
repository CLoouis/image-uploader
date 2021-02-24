package mongo

import (
	"context"

	"github.com/CLoouis/image-uploader/pkg/api/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	UserRepositoryImpl struct {
		userCollection *mongo.Collection
	}
)

func NewUserRepositoryImpl(userCollection *mongo.Collection) user.UserRepository {
	return &UserRepositoryImpl{userCollection: userCollection}
}

func (u *UserRepositoryImpl) Create(ctx context.Context, userData user.User) (string, error) {
	res, err := u.userCollection.InsertOne(ctx, userData)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (u *UserRepositoryImpl) FindById(ctx context.Context, userID string) (user.User, error) {
	var userData user.User
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return user.User{}, err
	}

	err = u.userCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&userData)
	if err != nil {
		return user.User{}, err
	}

	return userData, nil
}

func (u *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (user.User, error) {
	var userData user.User
	err := u.userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&userData)
	if err != nil {
		return user.User{}, err
	}

	return userData, nil
}

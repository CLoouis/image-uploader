package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	User struct {
		ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
		Email    string             `json:"email,omitempty" bson:"email,omitempty"`
		Password string             `json:"-"`
	}

	UserRepository interface {
		Create(context.Context, User) (string, error)
		FindById(context.Context, string) (User, error)
	}

	UserService interface {
		Create(context.Context, User) (string, error)
		Me(context.Context) (User, error)
	}
)

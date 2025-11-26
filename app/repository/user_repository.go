package repository

import (
	"context"
	"crud_alumni/app/model"
	"crud_alumni/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var FindUserByUsernameOrEmailFunc func(identifier string) (*model.User, string, error)

func FindUserByUsernameOrEmail(identifier string) (*model.User, string, error) {
	if FindUserByUsernameOrEmailFunc != nil {
		return FindUserByUsernameOrEmailFunc(identifier)
	}

	if database.UserCollection == nil {
		panic("❌ UserCollection belum diinisialisasi")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user model.User
	err := database.UserCollection.FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"username": identifier},
			{"email": identifier},
		},
	}).Decode(&user)
	if err != nil {
		return nil, "", err
	}
	return &user, user.PasswordHash, nil
}

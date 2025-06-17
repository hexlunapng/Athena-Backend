package models

import (
	db "Athena-Backend/database"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AccountID string             `bson:"accountId" json:"accountId"`
	Username  string             `bson:"username" json:"username"`
	DiscordID *string            `bson:"discordId,omitempty" json:"discordId,omitempty"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	Created   time.Time          `bson:"created" json:"created"`
	Banned    bool               `bson:"banned" json:"banned"`
}

func UserAccount(accountID, username, email, password string, discordID *string) *User {
	return &User{
		ID:        primitive.NewObjectID(),
		AccountID: accountID,
		Username:  username,
		DiscordID: discordID,
		Email:     email,
		Password:  password,
		Created:   time.Now(),
		Banned:    false,
	}
}

func (u *User) Save() error {
	collection := db.GetMongoCollection("users")
	_, err := collection.InsertOne(context.Background(), u)
	return err
}

func UserExists(email string) (bool, error) {
	collection := db.GetMongoCollection("users")
	filter := bson.M{"email": email}

	var user User
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

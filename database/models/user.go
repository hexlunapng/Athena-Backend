package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AccountID          string             `bson:"accountId" json:"accountId"`
	Username           string             `bson:"username" json:"username"`
	DiscordID          *string            `bson:"discordId,omitempty" json:"discordId,omitempty"`
	Email              string             `bson:"email" json:"email"`
	Password           string             `bson:"password" json:"-"`
	Created            time.Time          `bson:"created" json:"created"`
	Banned             bool               `bson:"banned" json:"banned"`
	LastUsernameChange *time.Time         `bson:"lastUsernameChange,omitempty" json:"lastUsernameChange,omitempty"`
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

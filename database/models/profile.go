package models

import (
	"context"
	"time"

	db "Athena-Backend/database"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Profiles struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	Created   time.Time              `bson:"created" json:"created"`
	Updated   *time.Time             `bson:"updated,omitempty" json:"updated,omitempty"`
	AccountID string                 `bson:"accountId" json:"accountId"`
	Profiles  map[string]interface{} `bson:"profiles" json:"profiles"`
}

func (p *Profiles) Save() error {
	collection := db.GetMongoCollection("profiles")
	_, err := collection.InsertOne(context.Background(), p)
	return err
}

func UserProfiles(accountID string, profiles map[string]interface{}) *Profiles {
	return &Profiles{
		ID:        primitive.NewObjectID(),
		Created:   time.Now(),
		AccountID: accountID,
		Profiles:  profiles,
	}
}

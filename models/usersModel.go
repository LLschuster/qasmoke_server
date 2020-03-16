package models

import (
	"context"
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// UserProfile represents the structure of a user profile document
type UserProfile struct {
	Name  string        `json:"name"`
	Genre []interface{} `json:"genre"`
}

// AddProfile appends a new user profile to the users collections
func AddProfile(profileCollection *mongo.Collection, data []byte) error {
	var user UserProfile
	err := json.Unmarshal(data, &user)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	profileCollection.InsertOne(ctx, user)
	return nil
}

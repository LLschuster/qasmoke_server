package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserProfile represents the structure of a user profile document
type UserProfile struct {
	UserID     string        `json:"userid"`
	Token      string        `json:"token"`
	Name       string        `json:"name"`
	Entitle    string        `json:"entitle"`
	About      string        `json:"about"`
	SampleLink string        `json:"sampleLink"`
	Genre      []interface{} `json:"genre"`
	Followers  []interface{} `json:"followers"`
	Comments   []interface{} `json:"comments"`
	Feed       []Post        `json:"feed"`
}

// AddProfile appends a new user profile to the users collections
func AddProfile(profileCollection *mongo.Collection, data []byte) error {
	var user UserProfile
	err := json.Unmarshal(data, &user)
	if err != nil {
		fmt.Println(err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	profileCollection.InsertOne(ctx, user)
	return nil
}

// GetUserProfile get a user profile base on a userID
func GetUserProfile(profileCollection *mongo.Collection, userID string) (*UserProfile, error) {
	var user UserProfile
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	filter := bson.M{"userid": userID}
	err := profileCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

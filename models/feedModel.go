package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

// PostRequest represents the structure of a PostRequest document
type PostRequest struct {
	Name     string        `json:"name"`
	Genre    []interface{} `json:"genre"`
	UserName string        `json:"userName"`
}

// Post represents the structure of the feed object in the database
type Post struct {
	Name   string        `json:"name"`
	Genre  []interface{} `json:"genre"`
	Author string        `json:"author"`
}

// GetUserFeeds gets the user feed of a user base on its userID
func GetUserFeeds(profileCollection *mongo.Collection, data []byte) ([]Post, error) {
	var (
		user      UserProfile
		foundUSer UserProfile
	)
	err := json.Unmarshal(data, &user)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	filter := bson.M{"userID": user.UserID}
	err = profileCollection.FindOne(ctx, filter).Decode(&foundUSer)
	if err != nil {
		return nil, err
	}
	return foundUSer.Feed, nil
}

// PublishNewPost Adds a new Post, and Updates the feed of every user that follows the same genre as the new Post.
func PublishNewPost(profileCollection *mongo.Collection, data []byte) error {
	var request PostRequest
	var usersToUpdate []*UserProfile
	err := json.Unmarshal(data, &request)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	usersToUpdate = getUserByGenres(ctx, profileCollection, request.Genre)
	var filter bson.M
	for _, user := range usersToUpdate {
		fmt.Printf("user to update %+v\n", user)
		filter = bson.M{"userID": user.UserID}
		update := bson.M{"$push": bson.M{"feed": bson.M{"name": request.Name, "author": request.UserName, "genre": request.Genre}}}
		result, err := profileCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("result fo update %+v\n", result)
	}
	return nil
}

func getUserByGenres(ctx context.Context, profileCollection *mongo.Collection, postGenre []interface{}) []*UserProfile {
	fmt.Printf("searching users in genre %v\n", postGenre)
	var userProfiles []*UserProfile
	var filter bson.M
	for _, value := range postGenre {
		filter = bson.M{"genre": value}
		cur, err := profileCollection.Find(ctx, filter)
		if err != nil {
			log.Fatal(err)
		}
		for cur.Next(ctx) {
			var user UserProfile
			cur.Decode(&user)
			userProfiles = append(userProfiles, &user)
		}
	}

	fmt.Printf("found user with profile %+v\n", userProfiles)
	return userProfiles
}

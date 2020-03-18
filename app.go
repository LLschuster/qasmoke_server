package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/llschuster/qasmoke/models"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
 ################## Login endpoints ###############
*/
func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Login endpoint"}`))
}

/*
 ################## user profile endpoints  ###############
*/
func createProfile(db *mongo.Database) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		users := db.Collection("users")
		params, err := ioutil.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "invalid body of request"}`))
		}

		err = models.AddProfile(users, params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Fatal error while Adding user to db"}`))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Profile was succesfully added"}`))
	}
	return fn
}

/*
 ################## user feed endpoints  ###############
*/
func postFeed(db *mongo.Database) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		users := db.Collection("users")
		params, err := ioutil.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "invalid body of request"}`))
		}

		err = models.PublishNewPost(users, params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Fatal error while Adding user to db"}`))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Profile was succesfully added"}`))
	}
	return fn
}

func getUserFeed(db *mongo.Database) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		users := db.Collection("users")
		params, err := ioutil.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "invalid body of request"}`))
		}

		feed, err := models.GetUserFeeds(users, params)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Fatal error while getting feed for user"}`))
		}

		w.WriteHeader(http.StatusOK)
		feedJSON, err := json.Marshal(feed)
		w.Write(feedJSON)
	}
	return fn
}

func main() {
	fmt.Println("http api")
	db, err := models.InitDb()
	if err != nil {
		fmt.Println(err)
	}
	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/login", login).Methods(http.MethodGet)

	api.HandleFunc("/profile", createProfile(db)).Methods(http.MethodPost)
	api.HandleFunc("/profile/feed", getUserFeed(db)).Methods(http.MethodGet)

	api.HandleFunc("/post", postFeed(db)).Methods(http.MethodPost)

	port := 5000
	fmt.Println("port ", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		fmt.Println(err)
	}
}

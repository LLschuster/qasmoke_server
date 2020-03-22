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
 ################## file upload endpoints ###############
*/
func uploadFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]
	fmt.Printf("user id %v\n", userID)

	r.ParseMultipartForm(5 << 20) // files are of max size 5 * 2^20 bytes ~= 5mb
	file, handler, err := r.FormFile("upload-file")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "Error by file upload"}`))
		return
	}
	defer file.Close()
	fmt.Printf("%v\n", handler.Filename)
	fmt.Printf("%v\n", handler.Header)
	fmt.Printf("%v\n", handler.Size)

	// Creates a tempfile, wish will be then added to s3 bucket. in the name the * will be replace by a random generated string
	// to delete the tempfile use defer os.Remove(fmt.sprintf("temp/%v",tempFile.Name()))
	tempfile, err := ioutil.TempFile("temp", "upload-*.png")
	defer tempfile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "Error by file upload"}`))
		return
	}
	tempfile.Write(fileBytes)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "file succesfully uploaded"}`))
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

	api.HandleFunc("/upload/{userID}", uploadFile).Methods(http.MethodPost)

	port := 5000
	fmt.Println("port ", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		log.Fatal(err)
	}
}

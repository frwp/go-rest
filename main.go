package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Note struct, main data in this app
type Note struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Details string `json:"details"`
	DueDate string `json:"dueDate"`
}

// Message struct, response sent to end user
type Message struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var db *gorm.DB

// load .env values
func loadDotEnv(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env")
	}

	return os.Getenv(key)
}

var DB_URL = loadDotEnv("DB_URL")
var DB_USER = loadDotEnv("DB_USER")
var DB_NAME = loadDotEnv("DB_NAME")

func dbConnect() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Jakarta", DB_URL, DB_USER, DB_NAME)
	_db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// Auto migrate
	_db.AutoMigrate(&Note{})

	if err != nil {
		log.Fatal(err)
	}

	return _db
}

func landing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello</h1>")
	fmt.Println("Landing hit")
}

func getPostNote(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case "GET":
		fmt.Println("Retreive notes")

		var notes []Note
		db.Find(&notes)
		w.Header().Add("content-type", "application/json")
		json.NewEncoder(w).Encode(notes)
	case "POST":
		reqBody, _ := ioutil.ReadAll(r.Body)
		var note Note

		json.Unmarshal(reqBody, &note)
		db.Create(&note)

		fmt.Println("Successfully created new note")
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(note)
	}
}

func getPutDeleteNote(w http.ResponseWriter, r *http.Request) {
	reqVars := mux.Vars(r)

	id, err := strconv.Atoi(reqVars["id"])
	w.Header().Set("content-type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		message := Message{
			Status:  "error",
			Message: "Bad Request in parameters",
		}
		json.NewEncoder(w).Encode(message)
		return
	}

	var note Note
	db.First(&note, id)

	if note.Id == 0 {
		w.WriteHeader(http.StatusNotFound)
		message := Message{
			Status:  "error",
			Message: "Notes not found!",
		}
		json.NewEncoder(w).Encode(message)
		return
	}

	switch method := r.Method; method {
	case "GET":
		fmt.Println("Find note by id")
		json.NewEncoder(w).Encode(note)
		return

	// Updates note with newNote
	case "PUT":
		fmt.Println("Update notes")
		var newNote Note
		reqBody, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(reqBody, &newNote)

		db.Model(&note).Updates(&newNote)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(newNote)
		return

	// Returns deleted note
	case "DELETE":
		fmt.Println("Delete notes")
		db.Delete(&note)
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(note)
	}
}

func handleRequest() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", landing)
	router.HandleFunc("/notes", getPostNote).Methods("GET", "POST")
	router.HandleFunc("/notes/{id}", getPutDeleteNote).Methods("GET", "PUT", "DELETE")
	log.Fatal(http.ListenAndServe("localhost:8000", router))
}

func main() {
	db = dbConnect()

	handleRequest()
}

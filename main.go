package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error

// Message structure
type Message struct {
	ID        uint      `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

// Connect to the database
func ConnectDB() {
	// Database connection
	connStr := "postgresql://neondb_owner:npg_54xTLPVuygmq@ep-wild-lab-a47ef400-pooler.us-east-1.aws.neon.tech/neondb?sslmode=require"
	db, err = gorm.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Connected to the database!")
}

// Initialize the tables
func InitTables() {
	// Create messages table if not exists
	db.AutoMigrate(&Message{})
}

// Middleware to handle CORS
func handleCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		// If it's a preflight OPTIONS request, return 200 OK
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Allow the actual request
		next.ServeHTTP(w, r)
	})
}

// Get all messages
func getMessages(w http.ResponseWriter, r *http.Request) {
	var messages []Message

	// Step 1: Get the last 10 messages in descending order
	if err := db.Order("created_at desc").Limit(10).Find(&messages).Error; err != nil {
		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
		return
	}

	// Step 2: Reverse the slice to show in ascending order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}




// Post a new message
func postMessage(w http.ResponseWriter, r *http.Request) {
	var message Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Save the message in the database
	if err := db.Create(&message).Error; err != nil {
		http.Error(w, "Failed to store message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func main() {
	ConnectDB()
	InitTables()

	// Set up the router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/messages", getMessages).Methods("GET")
	r.HandleFunc("/messages", postMessage).Methods("POST")

	// Apply the CORS middleware to the router
	http.Handle("/", handleCORS(r))

	// Start the server
	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

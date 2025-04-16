package main

import (
	"context"
	"encoding/json"
	"net/http"
)

func postMessage(w http.ResponseWriter, r *http.Request) {
	var m Message
	json.NewDecoder(r.Body).Decode(&m)

	err := DB.QueryRow(context.Background(),
		"INSERT INTO messages(text) VALUES($1) RETURNING id, created_at", m.Text).
		Scan(&m.ID, &m.CreatedAt)
	if err != nil {
		http.Error(w, "Failed to store message", 500)
		return
	}
	json.NewEncoder(w).Encode(m)
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query(context.Background(),
		"SELECT id, text, created_at FROM messages ORDER BY id DESC LIMIT 50")
	if err != nil {
		http.Error(w, "Failed to fetch messages", 500)
		return
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		rows.Scan(&m.ID, &m.Text, &m.CreatedAt)
		messages = append(messages, m)
	}
	json.NewEncoder(w).Encode(messages)
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/parnurzeal/gorequest"
)

type Message struct {
	ID      string
	Subject string `json:"subject"`
	Content string `json:"content"`
}

type Payload struct {
	Username  string `json:"username,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"`
	Channel   string `json:"channel,omitempty"`
	Text      string `json:"text,omitempty"`
}

func Send(webhookUrl string, proxy string, payload Payload) []error {
	resp, _, err := gorequest.New().
		Proxy(proxy).
		Post(webhookUrl).
		Send(payload).
		End()

	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return []error{fmt.Errorf("Error for sending a message. Status: %v", resp.Status)}
	}

	return nil
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {

	// Parse the request body to get the message data
	var m Message
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Failed to parse request body: %v", err)
		return
	}

	// Set an unique id for message
	guid := uuid.New()
	m.ID = guid.String()

	// Send message to slack
	webhookUrl := "your web hook url"

	payload := Payload{
		Text:      fmt.Sprintf("ID: %s\nSubject: %s\nContent: %s", m.ID, m.Subject, m.Content),
		Username:  "your username",
		Channel:   "#logs",
		IconEmoji: ":monkey_face:",
	}

	errs := Send(webhookUrl, "", payload)
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Printf("error: %s\n", err)
		}
	}

	// Return the created model as a response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(m)
}

func main() {
	// Routes and handlers
	http.HandleFunc("/tasks", CreateMessage)

	// Start the server
	log.Println("The server is listening on port 8000...")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("Server error: ", err)
	}
}

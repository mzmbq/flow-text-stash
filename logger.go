package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Error     string    `json:"error,omitempty"`
}

func StartServer(port string) error {
	http.HandleFunc("/logs", handleLogs)
	return http.ListenAndServe(":"+port, nil)
}

func handleLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var entry LogEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Here you could save logs to a file or database
	// For now, we'll just print them
	println("Log received:", entry.Timestamp.Format(time.RFC3339), entry.Message, entry.Error)
}

type Client struct {
	serverURL string
}

func NewClient(serverURL string) *Client {
	return &Client{serverURL: serverURL}
}

func (c *Client) LogError(err error) error {
	entry := LogEntry{
		Timestamp: time.Now(),
		Message:   "Error occurred",
		Error:     err.Error(),
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	_, err = http.Post(c.serverURL+"/logs", "application/json", bytes.NewBuffer(jsonData))
	return err
}

func (c *Client) LogMessage(msg string) error {
	entry := LogEntry{
		Timestamp: time.Now(),
		Message:   msg,
		Error:     "",
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	_, err = http.Post(c.serverURL+"/logs", "application/json", bytes.NewBuffer(jsonData))
	return err
}

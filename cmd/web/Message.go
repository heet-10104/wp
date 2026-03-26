package main

type MessageType int

const (
    Broadcast MessageType = iota 
    Personal                    
)

type BroadcastMessage int

const (
	Control BroadcastMessage = iota
	Chat
)

type Message struct {
	Sender  string          `json:"sender"`  // e.g., user ID or username
	Receiver string         `json:"receiver"` // e.g., user ID or username (for personal messages)
    Type    string          `json:"type"`    // e.g., "Personal", "broadcast"
    Payload string 			`json:"payload"` // The actual data
}
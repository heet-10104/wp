package main

import (
	"encoding/json"
	"log"
	"strings"
)

func (h *Hub) handleBroadcastMessage(msg Message) {

	for _, client := range h.clients {
		if strings.EqualFold(msg.Sender, client.id) {
			continue
		}
		payload, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error marshalling broadcast message: %v", err)
			continue
		}
		select {
		case client.send <- payload:
		default:
			close(client.send)
			delete(h.clients, client.id)
		}
	}
}
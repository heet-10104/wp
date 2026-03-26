//sanity check payload
// put it in db
package main

import (
	"encoding/json"
	"log"
)

func (h *Hub) handlePersonalMessage(msg Message) {
	log.Printf("Private message from %s to %s: %s", msg.Sender, msg.Receiver, msg.Payload)
	if target, ok := h.clients[msg.Receiver]; ok {
		payload, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error marshalling personal message: %v", err)
			return
		}
		select {
		case target.send <- payload:
		default:
			close(target.send)
			delete(h.clients, target.id)
		}
	}
}
package main

import "log"

func (h *Hub) registerConnection(client *Client) {
	//check if new new client is already in the db or not
	//if not add it to the db
	//if yes, update the db to show that client is online
	h.clients[client.id] = client
	log.Printf("User %s connected", client.id)
}

func (h *Hub) unregisterConnection(client *Client) {
	//update the db when client disconnects to show that client is offline
	if _, ok := h.clients[client.id]; ok {
		delete(h.clients, client.id)
		close(client.send)
	}
}
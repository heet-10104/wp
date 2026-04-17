package main

import "github.com/gorilla/websocket"
import "encoding/json"
import(
	"log"
	"net/http"
)

type Client struct {
	id   string
	conn *websocket.Conn
	send chan []byte
	room string
}

// get list of clients in a room end point is /clients?room=roomName

func getClients(w http.ResponseWriter, r *http.Request) {
	roomKeys, ok := r.URL.Query()["room"]
	if !ok || len(roomKeys[0]) < 1 {
		log.Println("Url Param 'room' is missing")
		return
	}
	room := roomKeys[0]

	if _, ok := rooms[room]; !ok {
		log.Println("Room does not exist")
		return
	}

	hub := rooms[room]
	clients := make([]string, 0)
	for _, client := range hub.clients {
		clients = append(clients, client.id)
	}
	clients = append(clients, "*")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clients)
}
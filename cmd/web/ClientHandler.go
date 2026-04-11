package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	roomParam := r.URL.Query()["room"]
	if len(roomParam) == 0 || len(roomParam[0]) < 1 {
		log.Println("Url Param 'room' is missing")
		return
	}
	room := roomParam[0]

	// Create room if it doesn't exist
	if _, ok := rooms[room]; !ok {
		rooms[room] = newHub()
		rooms[room].room = room
		go rooms[room].run()
		log.Printf("Creating new room: %s", room)
	}

	clients := make([]string, 0)
	for clientID := range rooms[room].clients {
		clients = append(clients, clientID)
	}

	data := struct {
		Room    string
		Clients []string
	}{
		Room:    room,
		Clients: clients,
	}

	log.Printf("Room %s has clients: %v", room, clients)

	// data := newTemplateData(r)
	render(w, http.StatusOK, "home.tmpl", data)
}

func joinRoom(w http.ResponseWriter, r *http.Request) {
	render(w, http.StatusOK, "join.tmpl", nil)
}

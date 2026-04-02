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

	//get room model by room name

	if _, ok := rooms[room]; !ok {
		log.Printf("Room %s does not exist", room)
		return
	}

	clients := make([]string, 0)
	for clientID := range rooms[room].clients {
		clients = append(clients, clientID)
	}

	data := struct{
		Room string
		Clients []string
	}{
		Room: room,
		Clients: clients,
	}

	log.Printf("Room %s has clients: %v", room, clients)

	// data := newTemplateData(r)
	render(w, http.StatusOK, "home.tmpl", data)
}
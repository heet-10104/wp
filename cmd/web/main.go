package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"github.com/gorilla/websocket"
)

/*BUGS
1. refreh page leads to empty msgs
2. information of the room must be pushed
3. need to implement seek, next, prev
*/

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var rooms = make(map[string]*Hub)

type Hub struct {
    clients    map[string]*Client 
	broadcast  chan Message
	private    chan Message
	register   chan *Client
	unregister chan *Client
	room       string
}

func newHub() *Hub {
	return &Hub{
        clients:    make(map[string]*Client),
		broadcast:  make(chan Message),
		private:    make(chan Message),
        register:   make(chan *Client),
        unregister: make(chan *Client),
	}
}

func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		if !validateMessage(msg) {
			log.Printf("Invalid message format: %+v", msg)
			continue
		}

		log.Printf("Received message: %+v", msg)

		switch msg.Type {
		case Personal:
			hub.private <- msg

		case Broadcast:
			hub.broadcast <- msg

		case Control:
			hub.broadcast <- msg

		default:
			log.Printf("Unknown message type after validation: %v", msg.Type)
		}
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			go h.registerConnection(client)

		case client := <-h.unregister:
			go h.unregisterConnection(client)

		case msg := <-h.private:
			go h.handlePersonalMessage(msg)

		case message := <-h.broadcast:
			go h.handleBroadcastMessage(message)

		}
	}
}

func serveWS(w http.ResponseWriter, r *http.Request) {
    // Get ID from URL: /ws?id=heet
    keys, ok := r.URL.Query()["username"]
	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'username' is missing")
		return
	}
	username := keys[0]

	roomKeys, ok := r.URL.Query()["room"]
	if !ok || len(roomKeys[0]) < 1 {
		log.Println("Url Param 'room' is missing")
		return
	}
	room := roomKeys[0]

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	if _, ok := rooms[room]; !ok {
		rooms[room] = newHub()
		rooms[room].room = room
		go rooms[room].run()
	}
	hub := rooms[room]

	client := &Client{
		id:   username,
		room: room,
		conn: conn,
		send: make(chan []byte, 256),
	}
	hub.register <- client

    go client.writePump()
    go client.readPump(hub)
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/ws", serveWS)
	log.Println("WebSocket server starting on :8080")
	go func() {
		err = http.ListenAndServe(":8080", nil)
		if err != nil {
			fmt.Println("Error starting server:", err)
		}
	}()

	addr := os.Getenv("ADDR")
	mux := routes()
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("Starting HTTP server on %s", addr)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatalf("Could not start HTTP server: %v", err)
		}
	}()

	select {}
}

func routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/home", home)
	router.HandlerFunc(http.MethodGet, "/", joinRoom)
	router.ServeFiles("/resources/*filepath", http.Dir("./resources"))

	return router
}
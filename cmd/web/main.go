package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"database/sql"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"github.com/gorilla/websocket"
)

//database initiilization
//server initialization
//ws initilization
//handler functions

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
       return true
    },
}

type Hub struct {
    clients    map[string]*Client 
    broadcast  chan Message
    private    chan Message 
    register   chan *Client
    unregister chan *Client
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
		log.Printf("Received message: %s", msg.Type)

		switch msg.Type {
		case "Personal":
			if msg.Receiver=="*"{
				log.Printf("Personal message must have a specific receiver, got: %s", msg.Receiver)
				continue
			}
			// Handle personal message (e.g., send to specific client)
			hub.private <- msg

		case "Broadcast":
			if msg.Receiver!="*"{
				log.Printf("Broadcast message should have receiver as '*', got: %s", msg.Receiver)
				continue
			}
			// Handle broadcast message (e.g., send to all clients)
			hub.broadcast <- msg

		default:
			log.Printf("Unknown message type: %s", msg.Type)
		}
		// Handle the messgae //async db write

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
			h.registerConnection(client)

        case client := <-h.unregister:
			h.unregisterConnection(client)

        case msg := <-h.private:
			h.handlePersonalMessage(msg)

        case message := <-h.broadcast:
			h.handleBroadcastMessage(message)

        }
    }
}

func serveWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
    // Get ID from URL: /ws?id=heet
    keys, ok := r.URL.Query()["id"]
    if !ok || len(keys[0]) < 1 {
        log.Println("Url Param 'id' is missing")
        return
    }
    userID := keys[0]

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }

    client := &Client{
        id:   userID,
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
		log.Fatal("Error loading .env file")
	}
	dsn := os.Getenv("DSN")
	log.Printf("Connecting to database with DSN: %s", dsn)
	_ = loadDatabase(dsn)
	log.Println("Database connected successfully")

	hub := newHub()
	go hub.run() // IMPORTANT

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(hub, w, r)
	})
    err = http.ListenAndServe(":8080", nil)
    if err != nil {
       fmt.Println("Error starting server:", err)
    }
	log.Println("Server started on :8080")
}

func loadDatabase(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}

	return db
}
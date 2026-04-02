package main

import "github.com/gorilla/websocket"

type Client struct {
	id   string
	conn *websocket.Conn
	send chan []byte
	room string
}
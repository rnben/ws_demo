package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[string]*websocket.Conn)
var broadcast = make(chan Message)

// Message struct for message passing
type Message struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}

var upgrader = websocket.Upgrader{}

func main() {
	http.HandleFunc("/", handleWSConnections)

	go handleBroadcast()

	log.Println("Server started on localhost:8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}

func handleWSConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading HTTP to WebSocket:", err)
		return
	}
	defer ws.Close()

	log.Println(ws.RemoteAddr(), "Connected to websocket")

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading message from WebSocket:", err)
			delete(clients, msg.Author)
			break
		}

		if msg.Content == "register" {
			clients[msg.Author] = ws
			log.Println(ws.RemoteAddr(), msg.Author, "registered successfully")
		} else if msg.Content == "ping" {
			broadcast <- msg
		} else {
			log.Println("invalid message")
			continue
		}
	}
}

// handleBroadcast sends messages to all connected clients
func handleBroadcast() {
	for {
		msg := <-broadcast

		client, ok := clients[msg.Author]
		if !ok {
			continue
		}

		err := client.WriteJSON(msg)
		if err != nil {
			log.Println("Error writing message to WebSocket:", err)
			client.Close()
			delete(clients, msg.Author)
		}
	}
}

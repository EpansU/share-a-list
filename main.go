package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Item struct {
	Name string
}

var (
	connections = []*websocket.Conn{}
	items       = []Item{}
)

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(handleConnection))
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{}
	upgrader.CheckOrigin = checkOrigin

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection", err)
		return
	}
	defer conn.Close()
	connections = append(connections, conn)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message: ", err)
			break
		}
		log.Println("Message type: ", messageType)
		log.Println("Message: ", p)
		item := Item{string(p)}
		items = append(items, item)
		jsonMsg, err := json.Marshal(items)
		if err != nil {
			log.Println("Error encoding message: ", err)
			break
		}
		for _, c := range connections {
			if err := c.WriteMessage(messageType, jsonMsg); err != nil {
				log.Println("Error writing message: ", err)
				break
			}
		}
	}
}

func checkOrigin(r *http.Request) bool {
	return true
}

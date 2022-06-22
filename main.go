package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
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

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message: ", err)
			break
		}
		log.Println("Message type: ", messageType)
		log.Println("Message: ", p)
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println("Error writing message: ", err)
			break
		}
	}
}

func checkOrigin(r *http.Request) bool {
	return true
}

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func sendHealthUpdates(conn *websocket.Conn) {
	for {
		_, err := http.Get("http://localhost:8080/api/health")
		status := "offline"
		if err == nil {
			status = "healthy"
		}
		if err := conn.WriteMessage(websocket.TextMessage, []byte(status)); err != nil {
			log.Println("write:", err)
			return
		}
		time.Sleep(5 * time.Second)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()
	sendHealthUpdates(conn)
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	log.Println("Server started at http://localhost:8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

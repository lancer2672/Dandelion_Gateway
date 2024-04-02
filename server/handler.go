package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("READ", err)
			return
		}
		log.Println("Message", messageType)
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println("WRITE", err)
			return
		}
	}
}
func forwardWebSocket() {

}

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type    int
	Message []byte
}

var list = make(chan Message)

func main() {
	http.HandleFunc("/", top)
	http.HandleFunc("/ws", handleWebSocket)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func top(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebScoket Upgrade Error: %v", err)
		return
	}

	go func() {
		for {
			t, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("ReadMessage Error: %v", err)
				return
			}
			list <- Message{Type: t, Message: msg}
		}
	}()

	go func() {
		for {
			message := <-list
			fmt.Println(message)
			if err := conn.WriteMessage(message.Type, message.Message); err != nil {
				log.Printf("WriteMessage Error: %v", err)
				conn.Close()
				return
			}
		}
	}()
}

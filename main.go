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

	intermediary := intermediary{msg: make(chan Message)}
	client := newClient(conn, intermediary)

	go client.read()
	go client.write()
}

type intermediary struct {
	msg chan Message
}

func (i *intermediary) run() {
	// クライアントAから送られてきたメッセージを読み取ってチャネルに入れる

	// チャネルに入っているメッセージを取り出してクライアントBに送る
}

type client struct {
	conn         *websocket.Conn
	intermediary intermediary
}

func newClient(
	conn *websocket.Conn,
	intermediary intermediary,
) *client {
	return &client{
		conn:         conn,
		intermediary: intermediary,
	}
}

func (c *client) read() {
	for {
		t, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("ReadMessage Error: %v", err)
			return
		}
		c.intermediary.msg <- Message{Type: t, Message: msg}
	}
}

func (c *client) write() {
	for {
		message := <- c.intermediary.msg
		fmt.Println(message)
		if err := c.conn.WriteMessage(message.Type, message.Message); err != nil {
			log.Printf("WriteMessage Error: %v", err)
			c.conn.Close()
			return
		}
	}
}

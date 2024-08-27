package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type    int
	Message []byte
}

func main() {
	intermediary := newIntermediary()
	go intermediary.run()

	http.HandleFunc("/", top)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(w, r, intermediary)
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func top(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func handleWebSocket(w http.ResponseWriter, r *http.Request, i *intermediary) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebScoket Upgrade Error: %v", err)
		return
	}

	client := newClient(conn, i)
	client.intermediary.register <- client

	go client.read()
	go client.write()
}

type intermediary struct {
	msg      chan Message
	register chan *client
	clients  map[*client]bool
}

func newIntermediary() *intermediary {
	return &intermediary{
		msg:      make(chan Message),
		register: make(chan *client),
		clients:  make(map[*client]bool),
	}
}

func (i *intermediary) run() {
	for {
		select {
		case client := <-i.register:
			i.clients[client] = true
		case message := <-i.msg:
			for client := range i.clients {
				select {
				case client.msg <- message:
				default:
					close(client.msg)
					delete(i.clients, client)
				}
			}
		}
	}
}

type client struct {
	conn         *websocket.Conn
	intermediary *intermediary
	msg          chan Message
}

func newClient(
	conn *websocket.Conn,
	intermediary *intermediary,
) *client {
	return &client{
		conn:         conn,
		intermediary: intermediary,
		msg:          make(chan Message),
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
		message := <-c.msg
		if err := c.conn.WriteMessage(message.Type, message.Message); err != nil {
			log.Printf("WriteMessage Error: %v", err)
			c.conn.Close()
			return
		}
	}
}

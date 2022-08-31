package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	pool := NewPool()
	go pool.Listen()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
		}

		fmt.Printf("Connection Established...\n")

		client := &Client{
			Conn: conn,
			Pool: pool,
		}

		pool.Register <- client

		client.Read()

	})

	log.Fatal(http.ListenAndServe(":9090", nil))
	fmt.Println("server running!")

}

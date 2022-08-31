package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var connections map[*websocket.Conn]bool = make(map[*websocket.Conn]bool)

var mu sync.Mutex

func main() {

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
		}

		fmt.Printf("%T\n", conn)
		fmt.Printf("Connection Established...\n")

		connections[conn] = true

		defer func() {
			mu.Lock()
			delete(connections, conn)
			mu.Unlock()

			conn.Close()

		}()

		// PROBLEM: possiblity of concurrent writing to
		// sockets produces ERROR, MUST FIX

		// LOOP LOGIC
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}

			log.Println(string(p))

			for c := range connections {
				if err := c.WriteMessage(messageType, p); err != nil {
					log.Println(err)
					return
				}
			}
		}
	})

	fmt.Println("server running!")
	log.Fatal(http.ListenAndServe(":9090", nil))

}

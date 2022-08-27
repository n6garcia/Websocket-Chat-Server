package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var connections map[*websocket.Conn]bool = make(map[*websocket.Conn]bool)

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

		defer conn.Close()

		defer delete(connections, conn)

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

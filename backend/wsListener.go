package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type wsListener struct {
	url string
}

func (w wsListener) listen(data chan []byte) {
	c, _, err := websocket.DefaultDialer.Dial(w.url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
			data <- message
		}
	}()
}

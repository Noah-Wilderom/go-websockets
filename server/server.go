package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Serve(pool *Pool, writer http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	worker := &Worker{pool: pool, conn: conn, payload: make(chan []byte, 256), job: new(Job)}
	worker.pool.register <- worker

	t := time.Now()
	fmt.Printf("[%02d:%02d:%02d] Serving websockets at 0.0.0.0:4001\n", t.Hour(), t.Minute(), t.Second())

	go worker.WritePump()
	go worker.ReadPump()
}

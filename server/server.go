package server

import (
	"github.com/gorilla/websocket"
	"github.com/oklog/ulid/v2"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool {
            // Replace "example.com" with your Laravel application's domain
            allowedOrigin := "http://localhost"
            origin := r.Header.Get("Origin")
            return origin == allowedOrigin
    	},
}

func Serve(pool *Pool, writer http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	worker := &Worker{id: ulid.Make().String(), pool: pool, conn: conn, payload: make(chan []byte, 256), job: new(Job)}
	worker.pool.register <- worker

	go worker.WritePump()
	go worker.ReadPump()
}

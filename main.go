package main

import (
	"flag"
	"github.com/Noah-Wilderom/go-websockets/server"
	"log"
	"net/http"
)

func main() {
	var addr = flag.String("addr", ":4001", "http service address")

	flag.Parse()

	pool := server.NewPool()

	go pool.Run()

	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		server.Serve(pool, writer, request)
	})

	err := http.ListenAndServe(*addr, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

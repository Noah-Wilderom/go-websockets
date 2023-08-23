package main

import (
	"flag"
	"fmt"
	"github.com/Noah-Wilderom/go-websockets/server"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	setLogging()

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

func setLogging() {
	t := time.Now()
	f, err := os.OpenFile(fmt.Sprintf("logs/%02d-%02d-%02d.log", t.Year(), t.Month(), t.Day()), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
}

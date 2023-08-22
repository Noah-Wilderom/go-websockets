package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type Worker struct {
	pool    *Pool
	conn    *websocket.Conn
	payload chan []byte
	job     *Job
}

func (worker *Worker) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		worker.conn.Close()
	}()
	for {
		select {
		case payload, ok := <-worker.payload:
			worker.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				worker.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			writer, err := worker.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			var job Job

			err = json.Unmarshal(payload, &job)
			if err != nil {
				log.Println(err)
			}

			writer.Write([]byte(job.id))

			n := len(worker.payload)
			for i := 0; i < n; i++ {
				writer.Write([]byte{'\n'})
				writer.Write(<-worker.payload)
			}

			if err := writer.Close(); err != nil {
				return
			}
		case <-ticker.C:
			worker.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := worker.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (worker *Worker) ReadPump() {
	defer func() {
		worker.pool.unregister <- worker
		err := worker.conn.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	worker.conn.SetReadLimit(maxMessageSize)
	worker.conn.SetReadDeadline(time.Now().Add(pongWait))
	worker.conn.SetPongHandler(func(string) error {
		worker.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, payload, err := worker.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error: %v", err)
			}
			break
		}

		var job Job

		err = json.Unmarshal(payload, &job)
		if err != nil {
			log.Println(err)
		}
		worker.pool.dispatch <- payload
	}
}

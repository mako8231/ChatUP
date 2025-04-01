package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func (svr *Server) HandleEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := svr.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Connection Error:", err.Error())
		return
	}

	svr.Mutex.Lock() //lock the thread to write
	svr.Clients[conn] = true
	svr.Mutex.Unlock()

	//configuring heartbeat
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	defer func() {
		log.Println("Closing connection...")

		svr.Mutex.Lock()
		delete(svr.Clients, conn)
		svr.Mutex.Unlock()

		//Clean closure
		conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(5*time.Second))
		conn.Close()

		log.Println("Succefully closed connection.")
	}()

	for {
		mt, message, err := conn.ReadMessage()

		//Close the connection message
		if err != nil {
			fmt.Println("Error While getting the message", err.Error())
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Read Error: ", err.Error())
			}
			break
		}

		if mt == websocket.CloseMessage {
			log.Println("Closure message received")
			break
		}

		go func(msg []byte) {
			defer func() {
				if r := recover(); r != nil {
					log.Println("Recovered: ", r)
				}
			}()
			data := svr.MessageHandler(message)
			err := svr.validateMessage(data)
			//if the data isnt valid:
			if err != nil {
				svr.WriteMessage([]byte(err.Error()))
			} else {
				text := data["username"].(string) + ": " + data["message"].(string)
				svr.WriteMessage([]byte("!msg" + text))
			}
		}(message)
	}

}

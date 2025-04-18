package server

import (
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

/*
*

	Server Struct

	Router: [HttpRouter],
	Clients: [Map containing websocket connected clients]
	MessageHandler: [function that handles the received message]
	Mutex: [Prevent deadlock]
	Upgrader: [websocket upgrader handling]

*
*/
type Server struct {
	Router         *mux.Router
	Clients        map[*websocket.Conn]bool
	MessageHandler func(message []byte) map[string]interface{}
	Mutex          *sync.RWMutex
	Upgrader       websocket.Upgrader
}

// Initialize the server and returns the memory address
func StartServer(MessageHandler func(messageBytes []byte) map[string]interface{}, port string) *Server {
	var svr Server

	svr.Router = &mux.Router{}
	svr.Clients = make(map[*websocket.Conn]bool)
	svr.MessageHandler = MessageHandler
	svr.Mutex = &sync.RWMutex{}
	svr.Upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		HandshakeTimeout: 10 * time.Second,
		ReadBufferSize:   4096,
		WriteBufferSize:  4096,
	}

	svr.Router.HandleFunc("/endpoint/", svr.HandleEndpoint).Methods("GET")

	fs := http.FileServer(http.Dir("./public/"))
	svr.Router.PathPrefix("/").Handler(fs)

	go func() {
		log.Println("Server is listening in http://127.0.0.1:" + port + "/")
		if err := http.ListenAndServe(":"+port, svr.Router); err != nil {
			log.Fatal("Failed to start the server: ", err.Error())
		}

	}()

	return &svr
}

func (server *Server) WriteMessage(message []byte) {
	server.Mutex.RLock()
	for conn := range server.Clients {
		conn.WriteMessage(websocket.TextMessage, message)
	}
	server.Mutex.RUnlock()
}

func (server *Server) validateMessage(msgData map[string]interface{}) error {
	username := msgData["username"]
	message := msgData["message"]

	if username == nil {
		return errors.New("invalid message formatting - missing field: Username")
	}

	if message == nil {
		return errors.New("invalid message formatting - missing field: Message")
	}

	return nil
}

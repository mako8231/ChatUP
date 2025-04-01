package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/mako8231/chatup/server"
)

// Message handler function
func messageHandler(messageBytes []byte) map[string]interface{} {
	var jsonMapData map[string]interface{}
	messageString := string(messageBytes)

	json.Unmarshal([]byte(messageString), &jsonMapData)
	fmt.Println(jsonMapData)

	return jsonMapData
}

func main() {
	//Run the main application with the following arguments:
	//./program <port>
	execArgs := os.Args
	port := "8080"

	fmt.Println(len(execArgs))

	if len(execArgs) >= 2 {
		port = execArgs[1]
	}

	srv := server.StartServer(messageHandler, port)

	// Send a ticker to check the broadcast lifetime
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		srv.Mutex.RLock()
		if len(srv.Clients) > 0 {
			srv.WriteMessage([]byte("Heartbeat: " + time.Now().String()))
		}
		srv.Mutex.RUnlock()
	}
}

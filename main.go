package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mako8231/chatup/server"
)

// Message handler function
func messageHandler(messageBytes []byte) {
	var jsonMapData map[string]interface{}
	messageString := string(messageBytes)

	json.Unmarshal([]byte(messageString), &jsonMapData)
	fmt.Println(jsonMapData)

}

func main() {
	srv := server.StartServer(messageHandler)

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

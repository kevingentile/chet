package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Message struct {
	// Type string `json:"type"`
	ID string `json:"id"`
}

func main() {
	router := gin.Default()
	var wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	wsUpgrader.CheckOrigin = func(r *http.Request) bool {
		return r.Host == "localhost:8000"
	}
	router.GET("/ws", func(c *gin.Context) {

		con, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}

		go func() {
			m, err := json.Marshal(Message{"ID"})
			if err != nil {
				log.Println(err)
			}
			for {
				if err := con.WriteMessage(websocket.TextMessage, m); err != nil {
					log.Println(err)
				}
				time.Sleep(2 * time.Second)
			}
		}()

		for {
			t, msg, err := con.ReadMessage()
			if err != nil {
				break
			}
			message := Message{}
			if err := json.Unmarshal(msg, &message); err != nil {
				log.Println(err)
			}
			fmt.Println(t, message)
		}
	})
	router.Run("localhost:8000")
}

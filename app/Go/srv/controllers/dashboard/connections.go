package dashboard

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Connection struct {
	Name   string
	Conn   *websocket.Conn
	Global *Dashboard
}

func (c *Connection) Read() {
	// Keep connection alive, wait for any request.
	for {
		if _, _, err := c.Conn.ReadMessage(); err != nil {
			zap.S().Info("Error on read message:", err.Error())
			break
		} else {
			c.Global.Results <- NewResult(
				"7d09c3f5-f50d-402b-9f4b-756030320264",
				"1",
				"John Wills",
				"2021-01-27 10:23:24")
		}
	}

	c.Conn.Close()
	c.Global.Leave <- c
}

func (c *Connection) Write(message *Result) {
	b, err := json.Marshal(message)
	if err != nil {
		zap.S().Fatal(err)
	}

	if err := c.Conn.WriteMessage(websocket.TextMessage, b); err != nil {
		zap.S().Info("Error on write message:", err.Error())
	}
}

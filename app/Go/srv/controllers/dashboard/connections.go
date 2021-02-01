package dashboard_controller

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
		}
	}

	c.Conn.Close()
	c.Global.Leave <- c
}

func (c *Connection) WriteResult(message *ResultMessage) {
	b, err := json.Marshal(message)
	if err != nil {
		zap.S().Fatal(err)
	}

	if err := c.Conn.WriteMessage(websocket.TextMessage, b); err != nil {
		zap.S().Info("Error on write message:", err.Error())
	}
}

func (c *Connection) WriteUnfinishedResult(message *UnfinishedResultMessage) {
	b, err := json.Marshal(message)
	if err != nil {
		zap.S().Fatal(err)
	}

	if err := c.Conn.WriteMessage(websocket.TextMessage, b); err != nil {
		zap.S().Info("Error on write message:", err.Error())
	}
}

func (c *Connection) WriteFinishedResult(message *FinishedResultMessage) {
	b, err := json.Marshal(message)
	if err != nil {
		zap.S().Fatal(err)
	}

	if err := c.Conn.WriteMessage(websocket.TextMessage, b); err != nil {
		zap.S().Info("Error on write message:", err.Error())
	}
}

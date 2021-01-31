package dashboard

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
)

type Dashboard struct {
	ConnHub map[string]*Connection
	Results chan *Result
	Join    chan *Connection
	Leave   chan *Connection
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  512,
	WriteBufferSize: 512,
	CheckOrigin: func(r *http.Request) bool {
		zap.S().Infof("%s %s%s %v", r.Method, r.Host, r.RequestURI, r.Proto)
		return r.Method == http.MethodGet
	},
}

func (d *Dashboard) ResultsHandler(w http.ResponseWriter, r *http.Request) {
	upgradedConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		zap.S().Info("Error on websocket connection:", err.Error())
		return
	}

	uuid, err := uuid.NewV4()
	if err != nil {
		zap.S().Fatal(err)
	}

	conn := &Connection{
		Name:   fmt.Sprintf("anon-%d", uuid),
		Conn:   upgradedConn,
		Global: d,
	}

	d.Join <- conn

	conn.Read()
}

func (d *Dashboard) Run() {
	for {
		select {
		case conn := <-d.Join:
			d.add(conn)
		case result := <-d.Results:
			d.broadcast(result)
		case conn := <-d.Leave:
			d.disconnect(conn)
		}
	}
}

func (d *Dashboard) add(conn *Connection) {
	if _, usr := d.ConnHub[conn.Name]; !usr {
		d.ConnHub[conn.Name] = conn
		zap.S().Info("%s joined the chat", conn.Name)
	}
}

func (d *Dashboard) disconnect(conn *Connection) {
	if _, usr := d.ConnHub[conn.Name]; usr {
		defer conn.Conn.Close()
		delete(d.ConnHub, conn.Name)
	}
}

func (d *Dashboard) broadcast(result *Result) {
	zap.S().Infof("Broadcast result: %s, %s, %s",
		result.SportsmenStartNumber,
		result.SportsmenName,
		result.Time)
	for _, conn := range d.ConnHub {
		conn.Write(result)
	}
}

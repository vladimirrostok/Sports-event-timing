package main

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Dashboard struct {
	connHub map[string]*Connection
	results chan *Result
	join    chan *Connection
	leave   chan *Connection
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  512,
	WriteBufferSize: 512,
	CheckOrigin: func(r *http.Request) bool {
		log.Printf("%s %s%s %v", r.Method, r.Host, r.RequestURI, r.Proto)
		return r.Method == http.MethodGet
	},
}

func (d *Dashboard) Handler(w http.ResponseWriter, r *http.Request) {
	upgradedConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("Error on websocket connection:", err.Error())
	}

	uuid, err := uuid.NewV4()
	if err != nil {
		log.Fatalln(err)
	}

	conn := &Connection{
		Name:   fmt.Sprintf("anon-%d", uuid),
		Conn:   upgradedConn,
		Global: d,
	}

	d.join <- conn

	conn.Read()
}

func (d *Dashboard) Run() {
	for {
		select {
		case conn := <-d.join:
			d.add(conn)
		case result := <-d.results:
			d.broadcast(result)
		case conn := <-d.leave:
			d.disconnect(conn)
		}
	}
}

func (d *Dashboard) add(conn *Connection) {
	if _, usr := d.connHub[conn.Name]; !usr {
		d.connHub[conn.Name] = conn
		log.Printf("%s joined the chat", conn.Name)
	}
}

func (d *Dashboard) disconnect(conn *Connection) {
	if _, usr := d.connHub[conn.Name]; usr {
		defer conn.Conn.Close()
		delete(d.connHub, conn.Name)
	}
}

func (d *Dashboard) broadcast(result *Result) {
	log.Printf("Broadcast result: %s, %s, %s, %s",
		result.SportsmenStartNumber,
		result.SportsmenName,
		result.Time)
	for _, conn := range d.connHub {
		conn.Write(result)
	}
}

func main() {
	log.Printf("Server listening on http://localhost%s", ":8080")

	c := &Dashboard{
		connHub: make(map[string]*Connection),
		results: make(chan *Result),
		join:    make(chan *Connection),
		leave:   make(chan *Connection),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Go Sports Events Timing!"))
	})

	http.HandleFunc("/dashboard", c.Handler)

	go c.Run()

	log.Fatalln(http.ListenAndServe(":8080", nil))
}

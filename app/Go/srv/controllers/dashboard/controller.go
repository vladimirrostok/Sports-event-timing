package dashboard_controller

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"net/http"
	"sports/backend/domain/models/result"
	"sports/backend/domain/models/sportsmen"
	"strconv"
)

type Dashboard struct {
	LastResults *[]ResultMessage
	ConnHub     map[string]*Connection
	Results     chan ResultMessage
	Finish      chan FinishMessage
	Join        chan *Connection
	Leave       chan *Connection
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

	for _, msg := range *d.LastResults {
		conn.Write(&msg)
	}

	d.Join <- conn

	conn.Read()
}

func (d *Dashboard) Run(db *gorm.DB) error {
	lastResults, err := result.GetLastTenResults(*db)
	if err != nil {
		return err
	}

	// Convert domain results into application level results.
	// Load results from DB on app startup.
	var resultsMessages []ResultMessage
	for _, result := range *lastResults {
		version := uint32(1)
		sportsmenFetched, err := sportsmen.GetSportsmen(*db, result.SportsmenID, &version)
		if err != nil {
			return err
		}

		msg := ResultMessage{
			ID:                   result.ID.String(),
			SportsmenStartNumber: strconv.Itoa(int(sportsmenFetched.StartNumber)),
			SportsmenName:        fmt.Sprintf("%s %s", sportsmenFetched.FirstName, sportsmenFetched.LastName),
			TimeStart:            result.TimeStart.String(),
		}
		resultsMessages = append(resultsMessages, msg)
	}

	d.LastResults = &resultsMessages

	for {
		select {
		case conn := <-d.Join:
			d.add(conn)
		case result := <-d.Results:
			d.broadcast(&result)
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

func (d *Dashboard) broadcast(result *ResultMessage) {
	// Update stored results to return latest data to recently joined customers.
	updatedResults := append(*d.LastResults, *result)
	d.LastResults = &updatedResults

	zap.S().Infof("Broadcast result: %s, %s, %s",
		result.SportsmenStartNumber,
		result.SportsmenName,
		result.TimeStart)
	for _, conn := range d.ConnHub {
		conn.Write(result)
	}
}

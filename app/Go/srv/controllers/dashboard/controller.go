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
)

type Dashboard struct {
	LastResults *[]ResultMessage
	ConnHub     map[string]*Connection
	Results     chan UnfinishedResultMessage
	Finish      chan FinishedResultMessage
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

	if *d.LastResults != nil {
		conn.WriteAllCurrentResults(d.LastResults)
	} else {
		conn.WriteResult(nil)
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

	// Serve stored results in an reverse order so that the latest result will come the last
	// the last result will be placed on top of table then.
	for _, result := range *lastResults {
		version := uint32(1)
		sportsmenFetched, err := sportsmen.GetSportsmen(*db, result.SportsmenID, &version)
		if err != nil {
			return err
		}

		msg := ResultMessage{
			ID:                   result.ID.String(),
			SportsmenStartNumber: sportsmenFetched.StartNumber,
			SportsmenName:        fmt.Sprintf("%s %s", sportsmenFetched.FirstName, sportsmenFetched.LastName),
			TimeStart:            result.TimeStart,
			TimeFinish:           nil,
		}

		if result.TimeFinish != nil {
			msg.TimeFinish = result.TimeFinish
		}

		resultsMessages = append(resultsMessages, msg)
	}

	d.LastResults = &resultsMessages

	for {
		select {
		case conn := <-d.Join:
			d.add(conn)
		case result := <-d.Results:
			d.broadcastResult(&result)
		case finish := <-d.Finish:
			d.broadcastFinish(&finish)
		case conn := <-d.Leave:
			d.disconnect(conn)
		}
	}
}

func (d *Dashboard) add(conn *Connection) {
	if _, usr := d.ConnHub[conn.Name]; !usr {
		d.ConnHub[conn.Name] = conn
		zap.S().Info("%s joined the dashboard", conn.Name)
	}
}

func (d *Dashboard) disconnect(conn *Connection) {
	if _, usr := d.ConnHub[conn.Name]; usr {
		defer conn.Conn.Close()
		delete(d.ConnHub, conn.Name)
	}
}

func (d *Dashboard) broadcastResult(result *UnfinishedResultMessage) {
	// Update stored results to return latest data to recently joined customers.
	resultMessage := ResultMessage{
		ID:                   result.ID,
		SportsmenStartNumber: result.SportsmenStartNumber,
		SportsmenName:        result.SportsmenName,
		TimeStart:            result.TimeStart,
	}
	updatedResults := append(*d.LastResults, resultMessage)
	d.LastResults = &updatedResults

	zap.S().Infof("Broadcast result: %s, %s, %s",
		result.SportsmenStartNumber,
		result.SportsmenName,
		result.TimeStart)
	for _, conn := range d.ConnHub {
		conn.WriteUnfinishedResult(result)
	}
}

func (d *Dashboard) broadcastFinish(finish *FinishedResultMessage) {
	// Update stored results to return latest data to recently joined customers.
	for index, result := range *d.LastResults {
		if result.ID == finish.ID {
			time := finish.TimeFinish
			(*d.LastResults)[index].TimeFinish = &time
		}
	}

	zap.S().Infof("Broadcast result: %s, %s, %s",
		finish.SportsmenStartNumber,
		finish.SportsmenName,
		finish.TimeFinish)
	for _, conn := range d.ConnHub {
		conn.WriteFinishedResult(finish)
	}
}

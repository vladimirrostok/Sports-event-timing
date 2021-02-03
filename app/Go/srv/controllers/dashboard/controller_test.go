package dashboard_controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sports/backend/domain/models/checkpoint"
	"sports/backend/domain/models/result"
	"sports/backend/domain/models/sportsmen"
	"sports/backend/srv/cmd/config"
	dashboard_controller "sports/backend/srv/controllers/dashboard"
	result_controller "sports/backend/srv/controllers/result"
	"sports/backend/srv/server"
	"sports/backend/srv/utils"
	"strings"
	"time"
)

var _ = Describe("Results controller", func() {
	// To change the flags on the default logger to show the code line for better understanding.
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Set up database connection using configuration details.
	absPath, _ := filepath.Abs("../../cmd/config/")
	cfg := config.Config{}
	viper.AddConfigPath(absPath)
	viper.SetConfigName("configuration")
	viper.ReadInConfig()
	viper.Unmarshal(&cfg)

	Describe("Results received on connecting to dashboard", func() {
		// Setup separate database connections so that in asynchronous tests run the transactions will not overlap each other.
		// Otherwise one test will commit db rollback while second test hasn't finished yet.
		conn, err := utils.GetDBConnection(
			cfg.DBDriver,
			cfg.DBUsername,
			cfg.DBPassword,
			cfg.DBPort,
			cfg.DBHost,
			cfg.DBName,
		)
		Expect(err).To(BeNil())

		// Set up the dashboard Websocket API module
		dashboard := &dashboard_controller.Dashboard{
			ConnHub: make(map[string]*dashboard_controller.Connection),
			Results: make(chan dashboard_controller.UnfinishedResultMessage),
			Finish:  make(chan dashboard_controller.FinishedResultMessage),
			Join:    make(chan *dashboard_controller.Connection),
			Leave:   make(chan *dashboard_controller.Connection),
		}

		srv := server.Server{}
		srv.Addr = cfg.APIAddress
		srv.DB = conn
		srv.Router = mux.NewRouter()
		srv.Dashboard = dashboard

		db := conn.Begin()
		srv.DB = db

		// Run the server when the database has been set up.
		go srv.Dashboard.Run(srv.DB)

		for srv.Dashboard.LastResults == nil {
			time.Sleep(1 * time.Second)
			log.Print("Waiting for the srv to load the data")
		}

		AfterEach(func() {
			_ = db.Rollback()
		})

		When("There are no results yet", func() {
			s := httptest.NewServer(http.HandlerFunc(srv.Dashboard.ResultsHandler))

			// Convert http://127.0.0.1 to ws://127.0.0.
			u := "ws" + strings.TrimPrefix(s.URL, "http")
			Specify("Empty array returned", func() {
				// Connect to the server
				ws, _, err := websocket.DefaultDialer.Dial(u, nil)
				Expect(err).To(BeNil())

				// Read response and check to see if it's what we expect.
				_, msg, err := ws.ReadMessage()
				Expect(err).To(BeNil())
				Expect(string(msg)).To(Equal("null"))

				s.Close()
				ws.Close()
			})
		})
	})

	Describe("Results received on connecting to dashboard", func() {
		conn, err := utils.GetDBConnection(
			cfg.DBDriver,
			cfg.DBUsername,
			cfg.DBPassword,
			cfg.DBPort,
			cfg.DBHost,
			cfg.DBName,
		)
		Expect(err).To(BeNil())

		// Set up the dashboard Websocket API module
		dashboard := &dashboard_controller.Dashboard{
			ConnHub: make(map[string]*dashboard_controller.Connection),
			Results: make(chan dashboard_controller.UnfinishedResultMessage),
			Finish:  make(chan dashboard_controller.FinishedResultMessage),
			Join:    make(chan *dashboard_controller.Connection),
			Leave:   make(chan *dashboard_controller.Connection),
		}

		srv := server.Server{}
		srv.Addr = cfg.APIAddress
		srv.DB = conn
		srv.Router = mux.NewRouter()
		srv.Dashboard = dashboard

		db := conn.Begin()
		srv.DB = db

		AfterEach(func() {
			_ = db.Rollback()
		})

		When("There are results already", func() {
			var unfinishedResult result.UnfinishedResult
			var finishedResult result.FinishedResult

			pendingCheckpoint := checkpoint.PendingCheckpoint{
				ID:   uuid.Must(uuid.NewV4()),
				Name: "Corridor1",
			}

			pendingSportsmen := sportsmen.PendingSportsmen{
				ID:          uuid.Must(uuid.NewV4()),
				FirstName:   "Vladimir",
				LastName:    "Andrianov",
				StartNumber: 101,
			}

			pendingSportsmen2 := sportsmen.PendingSportsmen{
				ID:          uuid.Must(uuid.NewV4()),
				FirstName:   "Name2",
				LastName:    "Lastname2",
				StartNumber: 102,
			}

			BeforeEach(func() {
				_, err := checkpoint.Create(*db, pendingCheckpoint)
				Expect(err).To(BeNil())

				_, err = sportsmen.Create(*db, pendingSportsmen)
				Expect(err).To(BeNil())

				pendingResult := result.PendingResult{
					ID:           uuid.Must(uuid.NewV4()),
					CheckpointID: pendingCheckpoint.ID,
					SportsmenID:  pendingSportsmen.ID,
					TimeStart:    time.Now().Unix(),
				}

				unfinishedResult.ID = pendingResult.ID
				unfinishedResult.SportsmenID = pendingResult.SportsmenID
				unfinishedResult.CheckpointID = pendingResult.CheckpointID
				unfinishedResult.TimeStart = pendingResult.TimeStart

				_, err = result.Create(*db, pendingResult)
				Expect(err).To(BeNil())

				_, err = sportsmen.Create(*db, pendingSportsmen2)
				Expect(err).To(BeNil())

				pendingResult2 := result.PendingResult{
					ID:           uuid.Must(uuid.NewV4()),
					CheckpointID: pendingCheckpoint.ID,
					SportsmenID:  pendingSportsmen2.ID,
					TimeStart:    time.Now().Unix(),
				}

				_, err = result.Create(*db, pendingResult2)
				Expect(err).To(BeNil())

				unfinishedResult := result.UnfinishedResult{
					ID:           pendingResult2.ID,
					SportsmenID:  pendingResult2.SportsmenID,
					CheckpointID: pendingResult2.CheckpointID,
					TimeStart:    pendingResult2.TimeStart,
					Version:      1,
				}

				timeFinish := time.Now().Unix()
				finishedResult.ID = pendingResult2.ID
				finishedResult.SportsmenID = pendingResult2.SportsmenID
				finishedResult.CheckpointID = pendingResult2.CheckpointID
				finishedResult.TimeStart = pendingResult2.TimeStart
				finishedResult.TimeFinish = &timeFinish

				_, err = result.AddFinishTime(*db, timeFinish, unfinishedResult)
				Expect(err).To(BeNil())

				// Start server after data has been created, that way server will load existing data on start.
				go srv.Dashboard.Run(srv.DB)

				for srv.Dashboard.LastResults == nil {
					time.Sleep(1 * time.Second)
					log.Print("Waiting for the srv to load the data")
				}
			})

			Specify("Results returned", func() {
				s := httptest.NewServer(http.HandlerFunc(srv.Dashboard.ResultsHandler))
				// Convert http://127.0.0.1 to ws://127.0.0.
				u := "ws" + strings.TrimPrefix(s.URL, "http")

				// Connect to the server
				ws, _, err := websocket.DefaultDialer.Dial(u, nil)
				Expect(err).To(BeNil())

				// Read response and check to see if it's what we expect.
				_, msg, err := ws.ReadMessage()
				Expect(err).To(BeNil())

				resultsReceived := []dashboard_controller.ResultMessage{}

				err = json.Unmarshal(msg, &resultsReceived)
				Expect(err).To(BeNil())

				// Make sure results and ORDER of results is correct too.
				Expect(resultsReceived).To(Equal([]dashboard_controller.ResultMessage{
					{
						ID:                   finishedResult.ID.String(),
						SportsmenStartNumber: pendingSportsmen2.StartNumber,
						SportsmenName:        fmt.Sprintf("%s %s", pendingSportsmen2.FirstName, pendingSportsmen2.LastName),
						TimeStart:            finishedResult.TimeStart,
						TimeFinish:           finishedResult.TimeFinish,
					},
					{
						ID:                   unfinishedResult.ID.String(),
						SportsmenStartNumber: pendingSportsmen.StartNumber,
						SportsmenName:        fmt.Sprintf("%s %s", pendingSportsmen.FirstName, pendingSportsmen.LastName),
						TimeStart:            unfinishedResult.TimeStart,
						TimeFinish:           nil,
					},
				}))

				s.Close()
				ws.Close()
			})
		})
	})

	Describe("Messages on new result and new finish added and the results state returned right after update, and results order", func() {
		conn, err := utils.GetDBConnection(
			cfg.DBDriver,
			cfg.DBUsername,
			cfg.DBPassword,
			cfg.DBPort,
			cfg.DBHost,
			cfg.DBName,
		)
		Expect(err).To(BeNil())

		// Set up the dashboard Websocket API module
		dashboard := &dashboard_controller.Dashboard{
			ConnHub: make(map[string]*dashboard_controller.Connection),
			Results: make(chan dashboard_controller.UnfinishedResultMessage),
			Finish:  make(chan dashboard_controller.FinishedResultMessage),
			Join:    make(chan *dashboard_controller.Connection),
			Leave:   make(chan *dashboard_controller.Connection),
		}

		srv := server.Server{}
		srv.Addr = cfg.APIAddress
		srv.DB = conn
		srv.Router = mux.NewRouter()
		srv.Dashboard = dashboard

		db := conn.Begin()
		srv.DB = db

		AfterEach(func() {
			_ = db.Rollback()
		})

		When("The connection is made", func() {
			var resultToFinish result.FinishedResult

			pendingCheckpoint := checkpoint.PendingCheckpoint{
				ID:   uuid.Must(uuid.NewV4()),
				Name: "Corridor1",
			}

			pendingSportsmen := sportsmen.PendingSportsmen{
				ID:          uuid.Must(uuid.NewV4()),
				FirstName:   "Vladimir",
				LastName:    "Andrianov",
				StartNumber: 101,
			}

			pendingSportsmen2 := sportsmen.PendingSportsmen{
				ID:          uuid.Must(uuid.NewV4()),
				FirstName:   "Name2",
				LastName:    "Lastname2",
				StartNumber: 102,
			}

			BeforeEach(func() {
				_, err := checkpoint.Create(*db, pendingCheckpoint)
				Expect(err).To(BeNil())

				_, err = sportsmen.Create(*db, pendingSportsmen)
				Expect(err).To(BeNil())

				_, err = sportsmen.Create(*db, pendingSportsmen2)
				Expect(err).To(BeNil())

				timeFinish := time.Now().Unix()
				resultToFinish.TimeFinish = &timeFinish

				// Start server after data has been created, that way server will load existing data on start.
				go srv.Dashboard.Run(srv.DB)

				for srv.Dashboard.LastResults == nil {
					time.Sleep(1 * time.Second)
					log.Print("Waiting for the srv to load the data")
				}
			})

			Specify("Messages, results returned", func() {
				s := httptest.NewServer(http.HandlerFunc(srv.Dashboard.ResultsHandler))
				// Convert http://127.0.0.1 to ws://127.0.0.
				u := "ws" + strings.TrimPrefix(s.URL, "http")

				// Connect to the server
				ws, _, err := websocket.DefaultDialer.Dial(u, nil)
				Expect(err).To(BeNil())

				// Read response and check to see if it's what we expect.
				// Ensure that state returned from server has no results.
				_, msg, err := ws.ReadMessage()
				Expect(err).To(BeNil())
				Expect(string(msg)).To(Equal("null"))

				// Test new result message.
				timeNow := time.Now().Unix()

				newReq := result_controller.NewResultRequest{
					CheckpointID: pendingCheckpoint.ID.String(),
					SportsmenID:  pendingSportsmen.ID.String(),
					Time:         timeNow,
				}

				requestBody, err := json.Marshal(newReq)
				Expect(err).To(BeNil())

				req, err := http.NewRequest("POST", "/results", bytes.NewBufferString(string(requestBody)))
				Expect(err).To(BeNil())

				rr := httptest.NewRecorder()
				handler := result_controller.AddResult(&srv)
				handler.ServeHTTP(rr, req)

				// Read response and check to see if it's what we expect.
				// Ensure that state returned from server has new result added.
				_, msg, err = ws.ReadMessage()
				Expect(err).To(BeNil())

				resultsReceived := dashboard_controller.ResultMessage{}

				err = json.Unmarshal(msg, &resultsReceived)
				Expect(err).To(BeNil())

				// Make sure results and ORDER of results is correct too.
				Expect(resultsReceived.SportsmenStartNumber).To(Equal(pendingSportsmen.StartNumber))
				Expect(resultsReceived.SportsmenName).To(Equal(fmt.Sprintf("%s %s", pendingSportsmen.FirstName, pendingSportsmen.LastName)))
				Expect(resultsReceived.TimeStart).To(Equal(newReq.Time))

				msgNewResultReceived := dashboard_controller.UnfinishedResultMessage{}
				err = json.Unmarshal(msg, &msgNewResultReceived)
				Expect(err).To(BeNil())

				// Make sure new result message is correct.
				Expect(msgNewResultReceived.SportsmenName).To(Equal(fmt.Sprintf("%s %s", pendingSportsmen.FirstName, pendingSportsmen.LastName)))
				Expect(msgNewResultReceived.SportsmenStartNumber).To(Equal(pendingSportsmen.StartNumber))
				Expect(msgNewResultReceived.TimeStart).To(Equal(timeNow))

				// Test finish message.
				// Add second unfinished result to have an result to finish.
				pendingResult2 := result.PendingResult{
					ID:           uuid.Must(uuid.NewV4()),
					CheckpointID: pendingCheckpoint.ID,
					SportsmenID:  pendingSportsmen2.ID,
					TimeStart:    time.Now().Unix(),
				}

				requestBody, err = json.Marshal(pendingResult2)
				Expect(err).To(BeNil())

				req, err = http.NewRequest("POST", "/results", bytes.NewBufferString(string(requestBody)))
				Expect(err).To(BeNil())

				rr = httptest.NewRecorder()
				handler = result_controller.AddResult(&srv)
				handler.ServeHTTP(rr, req)

				// Read received message after creating new result.
				// Add new result message was tested above, skip this iteration.
				_, msg, err = ws.ReadMessage()
				Expect(err).To(BeNil())

				// Finish the new created result
				resultToFinish.ID = pendingResult2.ID
				resultToFinish.SportsmenID = pendingResult2.SportsmenID
				resultToFinish.CheckpointID = pendingCheckpoint.ID
				resultToFinish.TimeStart = pendingResult2.TimeStart

				timeNow = time.Now().Unix()
				resultToFinish.TimeFinish = &timeNow

				finishReq := result_controller.FinishRequest{
					CheckpointID: pendingCheckpoint.ID.String(),
					SportsmenID:  resultToFinish.SportsmenID.String(),
					Time:         *resultToFinish.TimeFinish,
				}

				requestBody, err = json.Marshal(finishReq)
				Expect(err).To(BeNil())

				req, err = http.NewRequest("POST", "/finish", bytes.NewBufferString(string(requestBody)))
				Expect(err).To(BeNil())

				rr = httptest.NewRecorder()
				handler = result_controller.AddFinishTime(&srv)
				handler.ServeHTTP(rr, req)

				// Read new message after finishing result.
				// Read response and check to see if it's what we expect.
				// Ensure that state returned from server has new result added.
				_, msg, err = ws.ReadMessage()
				Expect(err).To(BeNil())

				msgNewFinishReceived := dashboard_controller.FinishedResultMessage{}
				err = json.Unmarshal(msg, &msgNewFinishReceived)
				Expect(err).To(BeNil())
				//
				// Make sure results and ORDER of results is correct too.
				Expect(msgNewFinishReceived.SportsmenName).To(Equal(fmt.Sprintf("%s %s", pendingSportsmen2.FirstName, pendingSportsmen2.LastName)))
				Expect(msgNewFinishReceived.SportsmenStartNumber).To(Equal(pendingSportsmen2.StartNumber))
				Expect(msgNewFinishReceived.TimeFinish).To(Equal(*resultToFinish.TimeFinish))

				// Read response and check to see if it's what we expect.
				// Ensure that state returned from server has new result added.
				ws, _, err = websocket.DefaultDialer.Dial(u, nil)
				Expect(err).To(BeNil())

				_, msg, err = ws.ReadMessage()
				Expect(err).To(BeNil())

				resultsArrReceived := []dashboard_controller.ResultMessage{}

				err = json.Unmarshal(msg, &resultsArrReceived)
				Expect(err).To(BeNil())

				// Make sure results and ORDER of results is correct too.
				Expect(resultsArrReceived[0].SportsmenStartNumber).To(Equal(pendingSportsmen.StartNumber))
				Expect(resultsArrReceived[0].SportsmenName).To(Equal(fmt.Sprintf("%s %s", pendingSportsmen.FirstName, pendingSportsmen.LastName)))
				Expect(resultsArrReceived[0].TimeStart).To(Equal(newReq.Time))

				Expect(resultsArrReceived[1].SportsmenStartNumber).To(Equal(pendingSportsmen2.StartNumber))
				Expect(resultsArrReceived[1].SportsmenName).To(Equal(fmt.Sprintf("%s %s", pendingSportsmen2.FirstName, pendingSportsmen2.LastName)))
				Expect(resultsArrReceived[1].TimeStart).To(Equal(resultToFinish.TimeStart))
				Expect(resultsArrReceived[1].TimeFinish).To(Equal(&newReq.Time))

				s.Close()
				ws.Close()
			})
		})
	})
})

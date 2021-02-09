package result_controller

import (
	"bytes"
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sports/backend/domain/models/checkpoint"
	"sports/backend/domain/models/result"
	"sports/backend/domain/models/sportsmen"
	"sports/backend/srv/cmd/config"
	dashboard_controller "sports/backend/srv/controllers/dashboard"
	"sports/backend/srv/server"
	"sports/backend/srv/utils"
)

var _ = Describe("Results controller", func() {
	var (
		db *gorm.DB
	)

	// Set up database connection using configuration details.
	absPath, _ := filepath.Abs("../../cmd/config/")
	cfg := config.Config{}
	viper.AddConfigPath(absPath)
	viper.SetConfigName("configuration")
	viper.ReadInConfig()
	viper.Unmarshal(&cfg)
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

	go srv.Dashboard.Run(srv.DB)

	BeforeEach(func() {
		db = conn.Begin()
		srv.DB = db
	})

	AfterEach(func() {
		_ = db.Rollback()
	})

	Describe("Creating new result", func() {
		When("New result request is sent", func() {
			var pendingResult result.PendingResult
			var pendingCheckpoint1 checkpoint.PendingCheckpoint
			var pendingSportsmen1 sportsmen.PendingSportsmen
			var pendingCheckpoint2 checkpoint.PendingCheckpoint
			var pendingSportsmen2 sportsmen.PendingSportsmen

			BeforeEach(func() {
				pendingCheckpoint1 = checkpoint.PendingCheckpoint{
					ID:   uuid.Must(uuid.NewV4()),
					Name: "Corridor1",
				}

				_, err := checkpoint.Create(*db, pendingCheckpoint1)
				Expect(err).To(BeNil())

				pendingSportsmen1 = sportsmen.PendingSportsmen{
					ID:          uuid.Must(uuid.NewV4()),
					FirstName:   "Vladimir",
					LastName:    "Andrianov",
					StartNumber: 101,
				}

				_, err = sportsmen.Create(*db, pendingSportsmen1)
				Expect(err).To(BeNil())

				pendingCheckpoint2 = checkpoint.PendingCheckpoint{
					ID:   uuid.Must(uuid.NewV4()),
					Name: "Corridor2",
				}

				_, err = checkpoint.Create(*db, pendingCheckpoint2)
				Expect(err).To(BeNil())

				pendingSportsmen2 = sportsmen.PendingSportsmen{
					ID:          uuid.Must(uuid.NewV4()),
					FirstName:   "Vladimir",
					LastName:    "Andrianov",
					StartNumber: 102,
				}

				_, err = sportsmen.Create(*db, pendingSportsmen2)
				Expect(err).To(BeNil())

				pendingResult = result.PendingResult{
					ID:           uuid.Must(uuid.NewV4()),
					CheckpointID: pendingCheckpoint1.ID,
					SportsmenID:  pendingSportsmen1.ID,
					TimeStart:    utils.MakeTimestampInMilliseconds(),
				}

				_, err = result.Create(*db,
					result.PendingResult{
						ID:           uuid.Must(uuid.NewV4()),
						CheckpointID: pendingCheckpoint1.ID,
						SportsmenID:  pendingSportsmen1.ID,
						TimeStart:    utils.MakeTimestampInMilliseconds(),
					})
				Expect(err).To(BeNil())
			})

			Specify("The response returned", func() {
				samples := []struct {
					CheckpointID string `json:"checkpoint_id"`
					SportsmenID  string `json:"sportsmen_id"`
					Time         int64  `json:"time_start"`
					statusCode   int
					errorMessage string
				}{
					{
						CheckpointID: pendingCheckpoint2.ID.String(),
						SportsmenID:  pendingSportsmen2.ID.String(),
						Time:         pendingResult.TimeStart,
						statusCode:   http.StatusOK,
						errorMessage: "",
					},
					{
						CheckpointID: uuid.Must(uuid.NewV4()).String(),
						SportsmenID:  pendingSportsmen2.ID.String(),
						Time:         pendingResult.TimeStart,
						statusCode:   http.StatusUnprocessableEntity,
						errorMessage: "Checkpoint does not exist",
					},
					{
						CheckpointID: pendingCheckpoint2.ID.String(),
						SportsmenID:  uuid.Must(uuid.NewV4()).String(),
						Time:         pendingResult.TimeStart,
						statusCode:   http.StatusUnprocessableEntity,
						errorMessage: "Sportsmen does not exist",
					},
					{
						CheckpointID: "",
						SportsmenID:  pendingSportsmen2.ID.String(),
						Time:         pendingResult.TimeStart,
						statusCode:   http.StatusUnprocessableEntity,
						errorMessage: "checkpoint_id: cannot be blank.",
					},
					{
						CheckpointID: pendingCheckpoint2.ID.String(),
						SportsmenID:  "",
						Time:         pendingResult.TimeStart,
						statusCode:   http.StatusUnprocessableEntity,
						errorMessage: "sportsmen_id: cannot be blank.",
					},
					{
						CheckpointID: pendingCheckpoint2.ID.String(),
						SportsmenID:  pendingSportsmen2.ID.String(),
						statusCode:   http.StatusUnprocessableEntity,
						errorMessage: "time_start: cannot be blank.",
					},
					{
						CheckpointID: pendingCheckpoint1.ID.String(),
						SportsmenID:  pendingResult.SportsmenID.String(),
						Time:         pendingResult.TimeStart,
						statusCode:   http.StatusUnprocessableEntity,
						errorMessage: "Result already exists",
					},
				}

				for _, s := range samples {
					newReq := NewResultRequest{
						CheckpointID: s.CheckpointID,
						SportsmenID:  s.SportsmenID,
						Time:         s.Time,
					}

					requestBody, err := json.Marshal(newReq)
					Expect(err).To(gomega.BeNil())

					req, err := http.NewRequest("POST", "/results", bytes.NewBufferString(string(requestBody)))
					Expect(err).To(gomega.BeNil())

					rr := httptest.NewRecorder()
					handler := AddResult(&srv)
					handler.ServeHTTP(rr, req)

					responseMap := make(map[string]interface{})

					err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
					Expect(err).To(gomega.BeNil())

					Expect(rr.Code).To(Equal(s.statusCode))

					if rr.Code == 200 {
						Expect(rr.Body.String()).ToNot(Equal(""))
					}

					if rr.Code != 200 {
						Expect(responseMap["error"]).To(Equal(s.errorMessage))
					}
				}
			})
		})
	})

	Describe("Updating unfinished result with finish time", func() {
		When("Finish request is sent", func() {
			var pendingResult result.PendingResult
			var pendingCheckpoint checkpoint.PendingCheckpoint
			var pendingSportsmen sportsmen.PendingSportsmen

			BeforeEach(func() {
				pendingCheckpoint = checkpoint.PendingCheckpoint{
					ID:   uuid.Must(uuid.NewV4()),
					Name: "Corridor1",
				}

				_, err := checkpoint.Create(*db, pendingCheckpoint)
				Expect(err).To(BeNil())

				pendingSportsmen = sportsmen.PendingSportsmen{
					ID:          uuid.Must(uuid.NewV4()),
					FirstName:   "Vladimir",
					LastName:    "Andrianov",
					StartNumber: 101,
				}

				_, err = sportsmen.Create(*db, pendingSportsmen)
				Expect(err).To(BeNil())

				pendingResult = result.PendingResult{
					ID:           uuid.Must(uuid.NewV4()),
					CheckpointID: pendingCheckpoint.ID,
					SportsmenID:  pendingSportsmen.ID,
					TimeStart:    utils.MakeTimestampInMilliseconds(),
				}

				_, err = result.Create(*db,
					result.PendingResult{
						ID:           uuid.Must(uuid.NewV4()),
						CheckpointID: pendingCheckpoint.ID,
						SportsmenID:  pendingSportsmen.ID,
						TimeStart:    utils.MakeTimestampInMilliseconds(),
					})
				Expect(err).To(BeNil())
			})

			Specify("The response returned", func() {
				samples := []struct {
					CheckpointID string `json:"checkpoint_id"`
					SportsmenID  string `json:"sportsmen_id"`
					Time         int64  `json:"time_finish"`
					statusCode   int
					errorMessage string
				}{
					{
						CheckpointID: pendingCheckpoint.ID.String(),
						SportsmenID:  pendingSportsmen.ID.String(),
						Time:         pendingResult.TimeStart,
						statusCode:   http.StatusOK,
						errorMessage: "",
					},
					{
						CheckpointID: uuid.Must(uuid.NewV4()).String(),
						SportsmenID:  pendingSportsmen.ID.String(),
						Time:         pendingResult.TimeStart,
						statusCode:   http.StatusUnprocessableEntity,
						errorMessage: "Result not found: Result does not exist",
					},
					{
						CheckpointID: pendingCheckpoint.ID.String(),
						SportsmenID:  uuid.Must(uuid.NewV4()).String(),
						Time:         pendingResult.TimeStart,
						statusCode:   http.StatusUnprocessableEntity,
						errorMessage: "Result not found: Result does not exist",
					},
					{
						CheckpointID: "",
						SportsmenID:  pendingSportsmen.ID.String(),
						Time:         pendingResult.TimeStart,
						statusCode:   http.StatusUnprocessableEntity,
						errorMessage: "checkpoint_id: cannot be blank.",
					},
					{
						CheckpointID: pendingCheckpoint.ID.String(),
						SportsmenID:  "",
						Time:         pendingResult.TimeStart,
						statusCode:   http.StatusUnprocessableEntity,
						errorMessage: "sportsmen_id: cannot be blank.",
					},
					{
						CheckpointID: pendingCheckpoint.ID.String(),
						SportsmenID:  pendingSportsmen.ID.String(),
						statusCode:   http.StatusUnprocessableEntity,
						errorMessage: "time_finish: cannot be blank.",
					},
					{
						CheckpointID: pendingCheckpoint.ID.String(),
						SportsmenID:  pendingResult.SportsmenID.String(),
						Time:         pendingResult.TimeStart,
						statusCode:   http.StatusUnprocessableEntity,
						errorMessage: "Result has finish time already",
					},
				}

				for _, s := range samples {
					newReq := FinishRequest{
						CheckpointID: s.CheckpointID,
						SportsmenID:  s.SportsmenID,
						Time:         s.Time,
					}

					requestBody, err := json.Marshal(newReq)
					Expect(err).To(gomega.BeNil())

					req, err := http.NewRequest("POST", "/results", bytes.NewBufferString(string(requestBody)))
					Expect(err).To(gomega.BeNil())

					rr := httptest.NewRecorder()
					handler := AddFinishTime(&srv)
					handler.ServeHTTP(rr, req)

					responseMap := make(map[string]interface{})

					err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
					Expect(err).To(gomega.BeNil())

					Expect(rr.Code).To(Equal(s.statusCode))

					if rr.Code == 200 {
						Expect(rr.Body.String()).ToNot(Equal(""))
					}

					if rr.Code != 200 {
						Expect(responseMap["error"]).To(Equal(s.errorMessage))
					}
				}
			})
		})
	})
})

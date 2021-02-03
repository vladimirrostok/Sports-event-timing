package sportsmen_controller

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sports/backend/srv/cmd/config"
	"sports/backend/srv/server"
	"sports/backend/srv/utils"
)

var _ = Describe("Sportsmens controller", func() {
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

	srv := server.Server{}
	srv.Addr = cfg.APIAddress
	srv.DB = conn
	srv.Router = mux.NewRouter()

	BeforeEach(func() {
		db = conn.Begin()
	})

	AfterEach(func() {
		_ = db.Rollback()
	})

	Describe("Creating new sportsmen", func() {
		sampleData := NewSportsmenRequest{
			StartNumber: 101,
			FirstName:   "Vladimir",
			LastName:    "Andrianov",
		}

		When("New sportsmen request is sent", func() {
			Specify("The response returned", func() {
				samples := []struct {
					firstName    string
					lastName     string
					startNumber  uint32
					statusCode   int
					errorMessage string
				}{
					{
						firstName:    sampleData.FirstName,
						lastName:     sampleData.LastName,
						startNumber:  sampleData.StartNumber,
						statusCode:   http.StatusOK,
						errorMessage: "",
					},
					{
						firstName:    "",
						lastName:     sampleData.LastName,
						startNumber:  sampleData.StartNumber,
						statusCode:   http.StatusUnprocessableEntity,
						errorMessage: "first_name: cannot be blank.",
					},
					{
						firstName:    sampleData.FirstName,
						lastName:     "",
						startNumber:  sampleData.StartNumber,
						statusCode:   http.StatusUnprocessableEntity,
						errorMessage: "last_name: cannot be blank.",
					},
					{
						firstName:    sampleData.FirstName,
						lastName:     sampleData.LastName,
						statusCode:   http.StatusUnprocessableEntity,
						errorMessage: "start_number: cannot be blank.",
					},
				}

				for _, s := range samples {
					loginRequest := NewSportsmenRequest{
						StartNumber: s.startNumber,
						FirstName:   s.firstName,
						LastName:    s.lastName,
					}

					requestBody, err := json.Marshal(loginRequest)
					Expect(err).To(gomega.BeNil())

					req, err := http.NewRequest("POST", "/sportsmens", bytes.NewBufferString(string(requestBody)))
					Expect(err).To(gomega.BeNil())

					rr := httptest.NewRecorder()
					handler := AddSportsmen(&srv)
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

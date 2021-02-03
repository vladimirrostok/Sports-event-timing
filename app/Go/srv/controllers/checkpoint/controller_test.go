package checkpoint_controller

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

var _ = Describe("Checkpoints controller", func() {
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

	Describe("Creating new checkpoint", func() {
		sampleData := NewCheckpointRequest{
			Name: "corridor1",
		}

		When("New checkpoint request is sent", func() {
			Specify("The response returned", func() {
				samples := []struct {
					name         string
					statusCode   int
					errorMessage string
				}{
					{
						name:         sampleData.Name,
						statusCode:   http.StatusOK,
						errorMessage: "",
					},
					{
						name:         "",
						statusCode:   http.StatusUnprocessableEntity,
						errorMessage: "name: cannot be blank.",
					},
				}

				for _, s := range samples {
					loginRequest := NewCheckpointRequest{
						Name: s.name,
					}

					requestBody, err := json.Marshal(loginRequest)
					Expect(err).To(gomega.BeNil())

					req, err := http.NewRequest("POST", "/checkpoints", bytes.NewBufferString(string(requestBody)))
					Expect(err).To(gomega.BeNil())

					rr := httptest.NewRecorder()
					handler := AddCheckpoint(&srv)
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

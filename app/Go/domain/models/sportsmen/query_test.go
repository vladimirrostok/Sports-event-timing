package sportsmen_test

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"path/filepath"
	"sports/backend/domain/models/sportsmen"
	"sports/backend/srv/cmd/config"
	"sports/backend/srv/utils"
)

var _ = Describe("Managing sportsmens", func() {
	var (
		db *gorm.DB
	)

	// Set up database connection using configuration details.
	absPath, _ := filepath.Abs("../../../srv/cmd/config/")
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

	BeforeEach(func() {
		db = conn.Begin()
	})

	AfterEach(func() {
		_ = db.Rollback()
	})

	Describe("Fetching a sportsmen", func() {
		var sportsmenID uuid.UUID
		var sampleData = sportsmen.Sportsmen{
			FirstName:   "Vladimir",
			LastName:    "Andrianov",
			StartNumber: 101,
		}

		When("sportsmen exists", func() {
			BeforeEach(func() {
				sportsmenID = uuid.Must(uuid.NewV4())

				err := db.Create(&sportsmen.Sportsmen{
					ID:          sportsmenID,
					FirstName:   sampleData.FirstName,
					LastName:    sampleData.LastName,
					StartNumber: sampleData.StartNumber,
					Version:     1,
				}).Error

				Expect(err).To(BeNil())
			})

			Specify("the sportsmen returned", func() {
				v := uint32(1)
				fetched, err := sportsmen.GetSportsmen(*db, sportsmenID, &v)
				Expect(err).To(BeNil())

				Expect(fetched.ID).To(Equal(sportsmenID))
				Expect(fetched.FirstName).To(Equal(sampleData.FirstName))
				Expect(fetched.LastName).To(Equal(sampleData.LastName))
				Expect(fetched.StartNumber).To(Equal(sampleData.StartNumber))
				Expect(fetched.Version).To(Equal(uint32(1)))
			})
		})
	})
})

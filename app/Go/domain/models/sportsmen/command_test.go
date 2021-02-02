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

	Describe("Creating a new sportsmen", func() {
		var pendingSportsmen sportsmen.PendingSportsmen

		BeforeEach(func() {
			pendingSportsmen = sportsmen.PendingSportsmen{
				ID:          uuid.Must(uuid.NewV4()),
				FirstName:   "Vladimir",
				LastName:    "Andrianov",
				StartNumber: 101,
			}
		})

		When("the sportsmen is created", func() {
			Specify("the returned event", func() {
				event, err := sportsmen.Create(*db, pendingSportsmen)
				Expect(err).To(BeNil())

				Expect(event).To(Equal(&sportsmen.SportsmenCreatedEvent{
					SportsmenID: pendingSportsmen.ID.String(),
					FirstName:   pendingSportsmen.FirstName,
					LastName:    pendingSportsmen.LastName,
					StartNumber: pendingSportsmen.StartNumber,
					Version:     1,
				}))
			})

			Specify("the sportsmen is persisted in the database", func() {
				_, err := sportsmen.Create(*db, pendingSportsmen)
				Expect(err).To(BeNil())

				fetched := sportsmen.Sportsmen{}
				err = db.Model(&fetched).Where("id = ?", pendingSportsmen.ID).Take(&fetched).Error
				Expect(err).To(BeNil())

				Expect(fetched.ID).To(Equal(pendingSportsmen.ID))
				Expect(fetched.FirstName).To(Equal(pendingSportsmen.FirstName))
				Expect(fetched.LastName).To(Equal(pendingSportsmen.LastName))
				Expect(fetched.StartNumber).To(Equal(pendingSportsmen.StartNumber))
				Expect(fetched.Version).To(Equal(uint32(1)))
			})
		})
	})
})

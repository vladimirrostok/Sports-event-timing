package result_test

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"path/filepath"
	"sports/backend/domain/models/checkpoint"
	"sports/backend/domain/models/result"
	"sports/backend/domain/models/sportsmen"
	"sports/backend/srv/cmd/config"
	"sports/backend/srv/utils"
	"time"
)

var _ = Describe("Managing results", func() {
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

	Describe("Fetching an unfinished result", func() {
		var resultID uuid.UUID
		var sportsmenID uuid.UUID
		var checkpointID uuid.UUID

		sampleTime := time.Now().Add(1 * time.Second).Unix()
		var sampleData = result.Result{
			TimeStart:  time.Now().Unix(),
			TimeFinish: &sampleTime,
		}

		When("unfinished result exists", func() {
			BeforeEach(func() {
				checkpointID = uuid.Must(uuid.NewV4())

				err := db.Create(&checkpoint.Checkpoint{
					ID:      checkpointID,
					Name:    "Corridor1",
					Version: 1,
				}).Error

				Expect(err).To(BeNil())

				sportsmenID = uuid.Must(uuid.NewV4())

				err = db.Create(&sportsmen.Sportsmen{
					ID:          sportsmenID,
					FirstName:   "Vladimir",
					LastName:    "Andrianov",
					StartNumber: 101,
					Version:     1,
				}).Error

				Expect(err).To(BeNil())

				resultID = uuid.Must(uuid.NewV4())

				err = db.Create(&result.Result{
					ID:           resultID,
					TimeStart:    sampleData.TimeStart,
					CheckpointID: checkpointID,
					SportsmenID:  sportsmenID,
					Version:      1,
				}).Error

				Expect(err).To(BeNil())
			})

			Specify("the result returned", func() {
				v := uint32(1)
				fetched, err := result.GetUnfinishedResult(*db, checkpointID, sportsmenID, &v)
				Expect(err).To(BeNil())
				Expect(fetched.ID).To(Equal(resultID))
				Expect(fetched.CheckpointID).To(Equal(checkpointID))
				Expect(fetched.SportsmenID).To(Equal(sportsmenID))
				Expect(fetched.TimeStart).To(Equal(sampleData.TimeStart))
				Expect(fetched.Version).To(Equal(uint32(1)))
			})
		})
	})

	Describe("Fetching last results", func() {
		When("More than 10 results are stored", func() {
			BeforeEach(func() {
				for i := 0; i <= 20; i++ {
					checkpointID := uuid.Must(uuid.NewV4())

					err := db.Create(&checkpoint.Checkpoint{
						ID:      checkpointID,
						Name:    "Corridor1",
						Version: 1,
					}).Error

					Expect(err).To(BeNil())

					sportsmenID := uuid.Must(uuid.NewV4())

					err = db.Create(&sportsmen.Sportsmen{
						ID:          sportsmenID,
						FirstName:   "Vladimir",
						LastName:    "Andrianov",
						StartNumber: 101,
						Version:     1,
					}).Error

					Expect(err).To(BeNil())

					err = db.Create(&result.Result{
						ID:           uuid.Must(uuid.NewV4()),
						CheckpointID: checkpointID,
						SportsmenID:  sportsmenID,
						TimeStart:    int64(i),
						Version:      1,
					}).Error
					Expect(err).To(BeNil())
				}
			})

			Specify("Only 10 last results returned ordered by TimeStart", func() {
				fetched, err := result.GetLastTenResults(*db)
				Expect(err).To(BeNil())
				Expect(len(*fetched)).To(Equal(10))
				Expect((*fetched)[0].TimeStart).To(Equal(int64(20)))
				Expect((*fetched)[1].TimeStart).To(Equal(int64(19)))
				Expect((*fetched)[2].TimeStart).To(Equal(int64(18)))
				Expect((*fetched)[3].TimeStart).To(Equal(int64(17)))
				Expect((*fetched)[4].TimeStart).To(Equal(int64(16)))
				Expect((*fetched)[5].TimeStart).To(Equal(int64(15)))
				Expect((*fetched)[6].TimeStart).To(Equal(int64(14)))
				Expect((*fetched)[7].TimeStart).To(Equal(int64(13)))
				Expect((*fetched)[8].TimeStart).To(Equal(int64(12)))
				Expect((*fetched)[9].TimeStart).To(Equal(int64(11)))

			})
		})
	})
})

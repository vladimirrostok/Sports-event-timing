package result_test

import (
	"errors"
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

	Describe("Creating a new result", func() {
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
				TimeStart:    time.Now().Unix(),
			}
		})

		When("the result is created", func() {
			Specify("the returned event", func() {
				event, err := result.Create(*db, pendingResult)
				Expect(err).To(BeNil())

				Expect(event).To(Equal(&result.ResultCreatedEvent{
					ResultID:     pendingResult.ID.String(),
					SportsmenID:  pendingResult.SportsmenID.String(),
					CheckpointID: pendingResult.CheckpointID.String(),
					TimeStart:    pendingResult.TimeStart,
					Version:      1,
				}))
			})

			Specify("the result is persisted in the database", func() {
				_, err := result.Create(*db, pendingResult)
				Expect(err).To(BeNil())

				fetched := result.Result{}
				err = db.Model(&fetched).Where("id = ?", pendingResult.ID).Take(&fetched).Error
				Expect(err).To(BeNil())

				var timeFinish *int64

				Expect(fetched.ID).To(Equal(pendingResult.ID))
				Expect(fetched.SportsmenID).To(Equal(pendingResult.SportsmenID))
				Expect(fetched.CheckpointID).To(Equal(pendingResult.CheckpointID))
				Expect(fetched.TimeStart).To(Equal(pendingResult.TimeStart))
				Expect(fetched.TimeFinish).To(Equal(timeFinish))
				Expect(fetched.Version).To(Equal(uint32(1)))
			})
		})

		When("the result already exists", func() {
			Specify("the error returned is of AlreadyExists domain error type", func() {
				_, err := result.Create(*db, pendingResult)
				Expect(err).To(BeNil())

				_, err = result.Create(*db, pendingResult)
				Expect(errors.As(err, &result.AlreadyExists{})).To(BeTrue())
			})
		})
	})

	Describe("Appending a result with finish time", func() {
		var pendingResult result.PendingResult
		var unfinishedResult result.UnfinishedResult
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
				TimeStart:    time.Now().Unix(),
			}

			_, err = result.Create(*db, pendingResult)
			Expect(err).To(BeNil())

			resultFetched := result.Result{}
			err = db.Model(&result.Result{}).Where("id = ?", pendingResult.ID).Take(&resultFetched).Error
			Expect(err).To(BeNil())

			unfinishedResult = result.UnfinishedResult{
				ID:           resultFetched.ID,
				SportsmenID:  resultFetched.SportsmenID,
				CheckpointID: resultFetched.CheckpointID,
				TimeStart:    resultFetched.TimeStart,
				Version:      resultFetched.Version,
			}
		})

		When("the result is updated", func() {
			Specify("the returned event", func() {
				time := time.Now().Unix()
				event, err := result.AddFinishTime(*db, time, unfinishedResult)
				Expect(err).To(BeNil())

				Expect(event).To(Equal(&result.ResultFinishedEvent{
					ResultID:   pendingResult.ID.String(),
					TimeFinish: time,
					Version:    2,
				}))
			})

			Specify("the result is persisted in the database", func() {
				time := time.Now().Unix()
				_, err := result.AddFinishTime(*db, time, unfinishedResult)
				Expect(err).To(BeNil())

				fetched := result.Result{}
				err = db.Model(&fetched).Where("id = ?", pendingResult.ID).Take(&fetched).Error
				Expect(err).To(BeNil())

				Expect(fetched.ID).To(Equal(pendingResult.ID))
				Expect(fetched.SportsmenID).To(Equal(pendingResult.SportsmenID))
				Expect(fetched.CheckpointID).To(Equal(pendingResult.CheckpointID))
				Expect(fetched.TimeStart).To(Equal(pendingResult.TimeStart))
				Expect(fetched.TimeFinish).To(Equal(&time))
				Expect(fetched.Version).To(Equal(uint32(2)))
			})
		})

		When("the finished result already exists", func() {
			Specify("the error returned is of AlreadyFinished domain error type", func() {
				time := time.Now().Unix()
				_, err := result.AddFinishTime(*db, time, unfinishedResult)
				Expect(err).To(BeNil())

				_, err = result.AddFinishTime(*db, time, unfinishedResult)
				Expect(errors.As(err, &result.AlreadyFinished{})).To(BeTrue())
			})
		})
	})
})

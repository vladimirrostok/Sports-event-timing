package checkpoint_test

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"path/filepath"
	"sports/backend/domain/models/checkpoint"
	"sports/backend/srv/cmd/config"
	"sports/backend/srv/utils"
)

var _ = Describe("Managing checkpoints", func() {
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

	Describe("Creating a new checkpoint", func() {
		var pendingCheckpoint checkpoint.PendingCheckpoint

		BeforeEach(func() {
			pendingCheckpoint = checkpoint.PendingCheckpoint{
				ID:   uuid.Must(uuid.NewV4()),
				Name: "Corridor1",
			}
		})

		When("the checkpoint is created", func() {
			Specify("the returned event", func() {
				event, err := checkpoint.Create(*db, pendingCheckpoint)
				Expect(err).To(BeNil())

				Expect(event).To(Equal(&checkpoint.CheckpointCreatedEvent{
					CheckpointID: pendingCheckpoint.ID.String(),
					Name:         pendingCheckpoint.Name,
					Version:      1,
				}))
			})

			Specify("the checkpoint is persisted in the database", func() {
				_, err := checkpoint.Create(*db, pendingCheckpoint)
				Expect(err).To(BeNil())

				fetched := checkpoint.Checkpoint{}
				err = db.Model(&fetched).Where("id = ?", pendingCheckpoint.ID).Take(&fetched).Error
				Expect(err).To(BeNil())

				Expect(fetched.ID).To(Equal(pendingCheckpoint.ID))
				Expect(fetched.Name).To(Equal(pendingCheckpoint.Name))
				Expect(fetched.Version).To(Equal(uint32(1)))
			})
		})
	})
})

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

	Describe("Fetching an checkpoint", func() {
		var checkpointID uuid.UUID
		var sampleData = checkpoint.Checkpoint{
			Name: "Corridor1",
		}

		When("checkpoint exists", func() {
			BeforeEach(func() {
				checkpointID = uuid.Must(uuid.NewV4())

				err := db.Create(&checkpoint.Checkpoint{
					ID:      checkpointID,
					Name:    sampleData.Name,
					Version: 1,
				}).Error

				Expect(err).To(BeNil())
			})

			Specify("the checkpoint returned", func() {
				v := uint32(1)
				fetched, err := checkpoint.GetCheckpoint(*db, checkpointID, &v)
				Expect(err).To(BeNil())

				Expect(fetched.ID).To(Equal(checkpointID))
				Expect(fetched.Name).To(Equal(sampleData.Name))
				Expect(fetched.Version).To(Equal(uint32(1)))
			})
		})
	})
})

package utils

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go.uber.org/zap"
	"sports/backend/domain/models/checkpoint"
	"sports/backend/domain/models/result"
	"sports/backend/domain/models/sportsmen"
	"time"
)

// GetDBConnection with the given configuration details.
func GetDBConnection(driver, username, password, port, host, database string) (*gorm.DB, error) {
	var err error
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, username, database, password)
	db, err := gorm.Open(driver, DBURL)
	if err != nil {
		zap.S().Fatal("Cannot connect to %s database ", driver)
		zap.S().Fatal("Error: ", err)
	} else {
		zap.S().Infof("Connected to the %s database ", driver)
	}

	// Database migration
	db.AutoMigrate(
		&result.Result{},
		&checkpoint.Checkpoint{},
		&sportsmen.Sportsmen{},
	)

	db.Model(&result.Result{}).AddForeignKey("checkpoint_id", "checkpoints(id)", "RESTRICT", "RESTRICT")
	db.Model(&result.Result{}).AddForeignKey("sportsmen_id", "sportsmens(id)", "RESTRICT", "RESTRICT")

	return db, nil
}

// Get Time in milliseconds.
// Reference to example here https://gobyexample.com/epoch .
func MakeTimestampInMilliseconds() int64 {
	return (time.Now().UnixNano() / int64(time.Millisecond))
}

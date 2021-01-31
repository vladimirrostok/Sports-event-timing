package utils

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go.uber.org/zap"
)

// GetDBConnection with the given configuration details.
func GetDBConnection(driver, username, password, port, host, database string) (*gorm.DB, error) {
	var err error
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, username, database, password)
	DB, err := gorm.Open(driver, DBURL)
	if err != nil {
		zap.S().Fatal("Cannot connect to %s database ", driver)
		zap.S().Fatal("Error: ", err)
	} else {
		zap.S().Infof("Connected to the %s database ", driver)
	}

	return DB, nil
}

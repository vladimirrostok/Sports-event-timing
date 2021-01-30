package utils

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
)

// GetDBConnection with the given configuration details.
func GetDBConnection(driver, username, password, port, host, database string) (*gorm.DB, error) {
	var err error
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, username, database, password)
	DB, err := gorm.Open(driver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", driver)
		log.Fatalln("Error:", err)
	} else {
		fmt.Printf("Connected to the %s database ", driver)
	}

	return DB, nil
}

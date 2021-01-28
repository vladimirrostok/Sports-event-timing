package sportsmen

import (
	"github.com/gofrs/uuid"
	"time"
)

// Sportsmen represents a persistence model for the sportsmen entity.
type Sportsmen struct {
	ID          uuid.UUID `gorm:"primary_key" json:"id"`
	StartNumber int       `gorm:"not null" json:"start_number"`
	FirstName   string    `gorm:"not null" json:"first_name"`
	LastName    string    `gorm:"not null" json:"last_name"`
	CreatedAt   time.Time `gorm:"default:now();not null" json:"created_at"`
	Version     uint32    `gorm:"not null" json:"version"`
}

package sportsmen

import (
	"github.com/gofrs/uuid"
)

// Sportsmen represents a persistence model for the sportsmen entity.
type Sportsmen struct {
	ID          uuid.UUID `gorm:"primary_key" json:"id"`
	StartNumber uint32    `gorm:"not null" json:"start_number"`
	FirstName   string    `gorm:"not null" json:"first_name"`
	LastName    string    `gorm:"not null" json:"last_name"`
	CreatedAt   int64     `gorm:"default:extract(epoch from now());not null" json:"created_at"`
	Version     uint32    `gorm:"not null" json:"version"`
}

// PendingSportsmen represents an event result about to sign up for an event.
type PendingSportsmen struct {
	ID          uuid.UUID `gorm:"primary_key" json:"id"`
	StartNumber uint32    `gorm:"not null" json:"start_number"`
	FirstName   string    `gorm:"not null" json:"first_name"`
	LastName    string    `gorm:"not null" json:"last_name"`
}

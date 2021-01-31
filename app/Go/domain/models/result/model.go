package result

import (
	"github.com/gofrs/uuid"
	"time"
)

// Result represents a persistence model for the event result.
type Result struct {
	ID           uuid.UUID  `gorm:"primary_key" json:"id"`
	CheckpointID uuid.UUID  `gorm:"not null" json:"checkpoint_id"`
	SportsmenID  uuid.UUID  `gorm:"not null" json:"sportsmen_id"`
	TimeStart    *time.Time `gorm:"not null" json:"time_start"`
	TimeFinish   *time.Time `json:"time_finish"`
	CreatedAt    time.Time  `gorm:"default:now();not null" json:"created_at"`
	Version      uint32     `gorm:"not null" json:"version"`
}

// PendingResult represents an event result about to create.
type PendingResult struct {
	ID           uuid.UUID  `gorm:"primary_key" json:"id"`
	CheckpointID uuid.UUID  `gorm:"not null" json:"checkpoint_id"`
	SportsmenID  uuid.UUID  `gorm:"not null" json:"sportsmen_id"`
	TimeStart    *time.Time `gorm:"not null" json:"time_start"`
}

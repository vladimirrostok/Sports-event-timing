package result

import (
	"github.com/gofrs/uuid"
	"time"
)

// Result represents a persistence model for the event results.
type Result struct {
	ID           uuid.UUID `gorm:"primary_key" json:"id"`
	CheckpointID int       `gorm:"not null" json:"checkpoint_id"`
	SportsmenID  int       `gorm:"not null" json:"sportsmen_id"`
	EventStateID int       `gorm:"not null" json:"event_state_id"`
	Time         time.Time `gorm:"not null" json:"created_at"`
	CreatedAt    time.Time `gorm:"default:now();not null" json:"created_at"`
	Version      uint32    `gorm:"not null" json:"version"`
}

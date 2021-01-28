package eventstate

import (
	"github.com/gofrs/uuid"
	"time"
)

// Action represents a persistence model for the event state.
type EventState struct {
	ID        uuid.UUID `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `gorm:"default:now();not null" json:"created_at"`
	Version   uint32    `gorm:"not null" json:"version"`
}

package checkpoint

import (
	"github.com/gofrs/uuid"
)

// Checkpoint represents a persistence model for the checkpoint entity.
type Checkpoint struct {
	ID        uuid.UUID `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt int64     `gorm:"default:extract(epoch from now());not null" json:"created_at"`
	Version   uint32    `gorm:"not null" json:"version"`
}

// PendingCheckpoint represents a checkpoint about to create.
type PendingCheckpoint struct {
	ID   uuid.UUID `gorm:"primary_key" json:"id"`
	Name string    `gorm:"not null" json:"name"`
}

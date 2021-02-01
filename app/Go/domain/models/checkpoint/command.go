package checkpoint

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/jinzhu/gorm"
	"strings"
)

// Create a new checkpoint.
func Create(db gorm.DB, pendingCheckpoint PendingCheckpoint) (*CheckpointCreatedEvent, error) {
	pendingCheckpoint.Name = strings.TrimSpace(pendingCheckpoint.Name)
	pendingCheckpoint.Name = strings.Title(pendingCheckpoint.Name)

	if err := validation.ValidateStruct(
		&pendingCheckpoint,
		validation.Field(&pendingCheckpoint.ID, validation.Required, is.UUIDv4),
		validation.Field(&pendingCheckpoint.Name, validation.Required),
	); err != nil {
		return nil, err
	}

	newCheckpoint := Checkpoint{
		ID:      pendingCheckpoint.ID,
		Name:    pendingCheckpoint.Name,
		Version: 1,
	}

	if err := db.Create(&Checkpoint{
		ID:      newCheckpoint.ID,
		Name:    newCheckpoint.Name,
		Version: newCheckpoint.Version,
	}).Error; err != nil {
		return nil, err
	}

	return &CheckpointCreatedEvent{
		CheckpointID: newCheckpoint.ID.String(),
		Name:         newCheckpoint.Name,
		Version:      newCheckpoint.Version,
	}, nil
}

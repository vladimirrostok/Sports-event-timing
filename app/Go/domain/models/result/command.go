package result

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/jinzhu/gorm"
)

// Create a new result.
func Create(db gorm.DB, pendingResult PendingResult) (*ResultCreated, error) {
	if err := validation.ValidateStruct(
		&pendingResult,
		validation.Field(&pendingResult.ID, validation.Required, is.UUIDv4),
		validation.Field(&pendingResult.CheckpointID, validation.Required, is.UUIDv4),
		validation.Field(&pendingResult.SportsmenID, validation.Required, is.UUIDv4),
		validation.Field(&pendingResult.TimeStart, validation.Required),
	); err != nil {
		return nil, err
	}

	newResult := Result{
		ID:           pendingResult.ID,
		CheckpointID: pendingResult.CheckpointID,
		SportsmenID:  pendingResult.SportsmenID,
		TimeStart:    pendingResult.TimeStart,
		Version:      1,
	}

	if err := db.Create(&Result{
		ID:           newResult.ID,
		CheckpointID: newResult.CheckpointID,
		SportsmenID:  newResult.SportsmenID,
		TimeStart:    newResult.TimeStart,
		Version:      1,
	}).Error; err != nil {
		return nil, err
	}

	return &ResultCreated{
		ResultID:     newResult.ID.String(),
		CheckpointID: newResult.CheckpointID.String(),
		SportsmenID:  newResult.SportsmenID.String(),
		TimeStart:    newResult.TimeStart.String(),
		Version:      1,
	}, nil
}

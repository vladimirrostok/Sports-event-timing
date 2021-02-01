package result

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/jinzhu/gorm"
	domain_errors "sports/backend/domain/errors"
)

// Create a new result.
func Create(db gorm.DB, pendingResult PendingResult) (*ResultCreatedEvent, error) {
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

	return &ResultCreatedEvent{
		ResultID:     newResult.ID.String(),
		CheckpointID: newResult.CheckpointID.String(),
		SportsmenID:  newResult.SportsmenID.String(),
		TimeStart:    newResult.TimeStart.String(),
		Version:      1,
	}, nil
}

func AddFinishTime(db gorm.DB, finishTime string, unfinishedResult UnfinishedResult) (*ResultFinishedEvent, error) {
	// Update attributes with `struct`, will only update non-zero fields.
	// Update attributes with `map` instead.
	// https://gorm.io/docs/update.html#Updates-multiple-columns
	result := db.Model(&Result{}).
		Where("id = ? AND version = ?",
			unfinishedResult.ID,
			1,
		).Updates(map[string]interface{}{"time_finish": finishTime, "version": unfinishedResult.Version + 1})

	if result.Error != nil {
		return nil, fmt.Errorf("Error adding finish time to the result: %w", result.Error)
	} else if result.RowsAffected != 1 {
		return nil, fmt.Errorf("State conflict: %w", domain_errors.StateConflict{})
	}

	return &ResultFinishedEvent{
		ResultID:   unfinishedResult.ID.String(),
		TimeFinish: finishTime,
		Version:    unfinishedResult.Version + 1,
	}, nil
}

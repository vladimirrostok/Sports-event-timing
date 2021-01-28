package result

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	domain_errors "sports/backend/domain/errors"
)

// GetResult fetches a result.
func GetResult(db gorm.DB, pk uuid.UUID, version *uint32) (*Result, error) {
	var result Result

	err := db.Model(&result).Where("id = ?", pk).Take(&result).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, fmt.Errorf("Result not found: %w", NotFound{})
	} else if version != nil && result.Version != *version {
		return nil, fmt.Errorf("Invalid version tag: %w", domain_errors.InvalidVersion{})
	} else if err != nil {
		return nil, fmt.Errorf("Error loading result: %w", err)
	}

	return &Result{
		ID:           result.ID,
		CheckpointID: result.CheckpointID,
		SportsmenID:  result.SportsmenID,
		EventStateID: result.EventStateID,
		Time:         result.Time,
		CreatedAt:    result.CreatedAt,
		Version:      result.Version,
	}, nil
}

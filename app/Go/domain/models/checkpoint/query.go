package checkpoint

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	domain_errors "sports/backend/domain/errors"
)

// GetCheckpoint fetches a checkpoint.
func GetCheckpoint(db gorm.DB, pk uuid.UUID, version *uint32) (*Checkpoint, error) {
	var checkpoint Checkpoint

	err := db.Model(&checkpoint).Where("id = ?", pk).Take(&checkpoint).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, fmt.Errorf("Checkpoint not found: %w", NotFound{})
	} else if version != nil && checkpoint.Version != *version {
		return nil, fmt.Errorf("Invalid version tag: %w", domain_errors.InvalidVersion{})
	} else if err != nil {
		return nil, fmt.Errorf("Error loading checkpoint: %w", err)
	}

	return &Checkpoint{
		ID:        checkpoint.ID,
		Name:      checkpoint.Name,
		CreatedAt: checkpoint.CreatedAt,
		Version:   checkpoint.Version,
	}, nil
}

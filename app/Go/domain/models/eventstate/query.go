package eventstate

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	domain_errors "sports/backend/domain/errors"
)

// GetEventState fetches a event state.
func GetEventState(db gorm.DB, pk uuid.UUID, version *uint32) (*EventState, error) {
	var eventState EventState

	err := db.Model(&eventState).Where("id = ?", pk).Take(&eventState).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, fmt.Errorf("Event state not found: %w", NotFound{})
	} else if version != nil && eventState.Version != *version {
		return nil, fmt.Errorf("Invalid version tag: %w", domain_errors.InvalidVersion{})
	} else if err != nil {
		return nil, fmt.Errorf("Error loading event state: %w", err)
	}

	return &EventState{
		ID:        eventState.ID,
		Name:      eventState.Name,
		CreatedAt: eventState.CreatedAt,
		Version:   eventState.Version,
	}, nil
}

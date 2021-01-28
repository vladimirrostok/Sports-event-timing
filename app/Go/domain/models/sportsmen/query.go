package sportsmen

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	domain_errors "sports/backend/domain/errors"
)

// GetSportsmen fetches a sportsmen.
func GetSportsmen(db gorm.DB, pk uuid.UUID, version *uint32) (*Sportsmen, error) {
	var sportsmen Sportsmen

	err := db.Model(&sportsmen).Where("id = ?", pk).Take(&sportsmen).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, fmt.Errorf("Sportsmen not found: %w", NotFound{})
	} else if version != nil && sportsmen.Version != *version {
		return nil, fmt.Errorf("Invalid version tag: %w", domain_errors.InvalidVersion{})
	} else if err != nil {
		return nil, fmt.Errorf("Error loading sportsmen: %w", err)
	}

	return &Sportsmen{
		ID:          sportsmen.ID,
		StartNumber: sportsmen.StartNumber,
		FirstName:   sportsmen.FirstName,
		LastName:    sportsmen.LastName,
		CreatedAt:   sportsmen.CreatedAt,
		Version:     sportsmen.Version,
	}, nil
}

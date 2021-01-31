package result

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
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

	return &result, nil
}

func GetLastTenResults(db gorm.DB) (*[]Result, error) {
	var results []Result

	err := db.Order("time desc").Limit(10).Find(&results).Error
	if gorm.IsRecordNotFoundError(err) {
		return &results, nil
	} else if err != nil {
		zap.S().Info(err)
		return nil, fmt.Errorf("Error loading results: %w", err)
	}

	return &results, nil
}

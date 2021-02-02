package result

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	domain_errors "sports/backend/domain/errors"
)

// GetUnfinishedResult fetches a result.
func GetUnfinishedResult(db gorm.DB, checkpoint_id, sportsmen_id uuid.UUID, version *uint32) (*UnfinishedResult, error) {
	var result Result

	err := db.Model(&result).Where("checkpoint_id = ? AND sportsmen_id = ?", checkpoint_id, sportsmen_id).Take(&result).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, fmt.Errorf("Result not found: %w", NotFound{})
	} else if result.TimeFinish != nil {
		return nil, AlreadyFinished{}
	} else if version != nil && result.Version != *version {
		return nil, fmt.Errorf("Result has been already updated: %w", domain_errors.InvalidVersion{})
	} else if err != nil {
		return nil, fmt.Errorf("Error loading result: %w", err)
	}

	return &UnfinishedResult{
		ID:           result.ID,
		SportsmenID:  result.SportsmenID,
		CheckpointID: result.CheckpointID,
		TimeStart:    result.TimeStart,
		Version:      result.Version,
	}, nil
}

func GetLastTenResults(db gorm.DB) (*[]Result, error) {
	var results []Result

	err := db.Order("time_start desc").Limit(10).Find(&results).Error
	if gorm.IsRecordNotFoundError(err) {
		return &results, nil
	} else if err != nil {
		zap.S().Info(err)
		return nil, fmt.Errorf("Error loading results: %w", err)
	}

	return &results, nil
}

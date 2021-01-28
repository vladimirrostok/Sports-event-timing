package sportsmen

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/jinzhu/gorm"
)

// Create a new sportsmen.
func Create(db gorm.DB, pendingSportsmen PendingSportsmen) (*SportsmenCreated, error) {
	if err := validation.ValidateStruct(
		&pendingSportsmen,
		validation.Field(&pendingSportsmen.ID, validation.Required, is.UUIDv4),
		validation.Field(&pendingSportsmen.StartNumber, validation.Required),
		validation.Field(&pendingSportsmen.FirstName, validation.Required),
		validation.Field(&pendingSportsmen.LastName, validation.Required),
	); err != nil {
		return nil, err
	}

	newSportsmen := Sportsmen{
		ID:          pendingSportsmen.ID,
		StartNumber: pendingSportsmen.StartNumber,
		FirstName:   pendingSportsmen.FirstName,
		LastName:    pendingSportsmen.LastName,
		Version:     1,
	}

	if err := db.Create(&Sportsmen{
		ID:          newSportsmen.ID,
		StartNumber: newSportsmen.StartNumber,
		FirstName:   newSportsmen.FirstName,
		LastName:    newSportsmen.LastName,
		Version:     newSportsmen.Version,
	}).Error; err != nil {
		return nil, err
	}

	return &SportsmenCreated{
		SportsmenID: newSportsmen.ID.String(),
		StartNumber: newSportsmen.StartNumber,
		FirstName:   newSportsmen.FirstName,
		LastName:    newSportsmen.LastName,
		Version:     newSportsmen.Version,
	}, nil
}

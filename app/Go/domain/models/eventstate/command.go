package eventstate

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/jinzhu/gorm"
	"strings"
)

// Create a new event state.
func Create(db gorm.DB, pendingEventState PendingEventState) (*EventStateCreated, error) {
	pendingEventState.Name = strings.TrimSpace(pendingEventState.Name)
	pendingEventState.Name = strings.Title(pendingEventState.Name)

	if err := validation.ValidateStruct(
		&pendingEventState,
		validation.Field(&pendingEventState.ID, validation.Required, is.UUIDv4),
		validation.Field(&pendingEventState.Name, validation.Required),
	); err != nil {
		return nil, err
	}

	newEventState := EventState{
		ID:      pendingEventState.ID,
		Name:    pendingEventState.Name,
		Version: 1,
	}

	if err := db.Create(&EventState{
		ID:      newEventState.ID,
		Name:    newEventState.Name,
		Version: newEventState.Version,
	}).Error; err != nil {
		return nil, err
	}

	return &EventStateCreated{
		EventStateID: newEventState.ID.String(),
		Name:         newEventState.Name,
		Version:      newEventState.Version,
	}, nil
}

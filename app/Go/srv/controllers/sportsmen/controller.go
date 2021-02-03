package sportsmen_controller

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofrs/uuid"
	"io/ioutil"
	"net/http"
	"sports/backend/domain/models/sportsmen"
	"sports/backend/srv/responses"
	"sports/backend/srv/server"
)

// AddSportsmen handles the new sportsmen request.
func AddSportsmen(server *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		req := NewSportsmenRequest{}
		err = json.Unmarshal(body, &req)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		err = validation.ValidateStruct(&req,
			validation.Field(&req.StartNumber, validation.Required),
			validation.Field(&req.FirstName, validation.Required),
			validation.Field(&req.LastName, validation.Required),
		)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		newSportsmen := sportsmen.PendingSportsmen{
			ID:          uuid.Must(uuid.NewV4()),
			StartNumber: req.StartNumber,
			FirstName:   req.FirstName,
			LastName:    req.LastName,
		}

		sportsmenCreatedEvent, err := sportsmen.Create(*server.DB, newSportsmen)
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		responses.JSON(w, http.StatusOK, CreatedResponse{ID: sportsmenCreatedEvent.SportsmenID})
	}
}

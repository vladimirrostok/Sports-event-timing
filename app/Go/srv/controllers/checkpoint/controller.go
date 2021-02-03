package checkpoint_controller

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofrs/uuid"
	"io/ioutil"
	"net/http"
	"sports/backend/domain/models/checkpoint"
	"sports/backend/srv/responses"
	"sports/backend/srv/server"
)

// AddCheckpoint handles the new checkpoint request.
func AddCheckpoint(server *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		req := NewCheckpointRequest{}
		err = json.Unmarshal(body, &req)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		err = validation.ValidateStruct(&req,
			validation.Field(&req.Name, validation.Required),
		)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		newCheckpoint := checkpoint.PendingCheckpoint{
			ID:   uuid.Must(uuid.NewV4()),
			Name: req.Name,
		}

		checkpointCreatedEvent, err := checkpoint.Create(*server.DB, newCheckpoint)
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		responses.JSON(w, http.StatusOK, CreatedResponse{ID: checkpointCreatedEvent.CheckpointID})
	}
}

package result_controller

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gofrs/uuid"
	"io/ioutil"
	"net/http"
	"sports/backend/domain/models/result"
	"sports/backend/srv/controllers/dashboard"
	"sports/backend/srv/responses"
	"sports/backend/srv/server"
	"time"
)

// AddResult handles the new result.
func AddResult(server *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		req := NewResult{}
		err = json.Unmarshal(body, &req)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		err = validation.ValidateStruct(&req,
			validation.Field(&req.CheckpointID, validation.Required, is.UUIDv4),
			validation.Field(&req.SportsmenID, validation.Required, is.UUIDv4),
			validation.Field(&req.Time, validation.Required),
		)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		newResult := result.PendingResult{
			ID:           uuid.Must(uuid.NewV4()),
			CheckpointID: uuid.Must(uuid.FromString(req.CheckpointID)),
			SportsmenID:  uuid.Must(uuid.FromString(req.SportsmenID)),
			TimeStart:    nil,
		}

		if req.Time != "" {
			t, err := time.Parse(time.RFC3339, req.Time)
			newResult.TimeStart = &t
			if err != nil {
				responses.ERROR(w, http.StatusUnprocessableEntity, err)
				return
			}
		}

		_, err = result.Create(*server.DB, newResult)
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		server.Dashboard.Results <- dashboard_controller.ResultMessage{
			ID:                   newResult.ID.String(),
			SportsmenName:        newResult.ID.String(),
			SportsmenStartNumber: newResult.ID.String(),
			TimeStart:            newResult.TimeStart.String(),
		}

		responses.JSON(w, http.StatusOK, nil)
	}
}

// GetLastTenResults handles the latest results request.
func GetLastTenResults(server *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		results, err := result.GetLastTenResults(*server.DB)
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, nil)
			return
		}

		responses.JSON(w, http.StatusOK, results)
	}
}

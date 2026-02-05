package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dmitry/taskmanager/internal/dto"
	"github.com/dmitry/taskmanager/pkg/errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(dto.SuccessResponse{
		Success: true,
		Data:    data,
	})
}

func RespondError(w http.ResponseWriter, err error) {
	appErr, ok := err.(*errors.AppError)
	if !ok {
		appErr = errors.Internal(err, "Произошла непредвиденная ошибка")
	}

	response := dto.ErrorResponse{
		Success: false,
		Error: dto.ErrorDetailWrapper{
			Code:    string(appErr.Code),
			Message: appErr.Message,
			Details: appErr.Details,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.HTTPStatusCode())
	json.NewEncoder(w).Encode(response)
}

func ParseUUID(w http.ResponseWriter, r *http.Request, param string) (uuid.UUID, bool) {
	vars := mux.Vars(r)
	idStr := vars[param]

	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(w, errors.BadRequest("Неверный формат UUID"))
		return uuid.Nil, false
	}

	return id, true
}

func DecodeJSON(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		RespondError(w, errors.BadRequest("Неверный формат JSON"))
		return false
	}
	return true
}

package api

import (
	"encoding/json"
	"errors"
	"net/http"

	apperrors "github.com/mytheresa/go-hiring-challenge/errors"
)

func OKResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func CreatedResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}

func ErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func HandleError(w http.ResponseWriter, err error) {
	var appErr *apperrors.AppError
	if !errors.As(err, &appErr) {
		ErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	switch appErr.Kind {
	case apperrors.KindInvalidInput:
		ErrorResponse(w, http.StatusBadRequest, appErr.Message)
	case apperrors.KindNotFound:
		ErrorResponse(w, http.StatusNotFound, appErr.Message)
	case apperrors.KindAlreadyExists:
		ErrorResponse(w, http.StatusConflict, appErr.Message)
	default:
		ErrorResponse(w, http.StatusInternalServerError, appErr.Message)
	}
}

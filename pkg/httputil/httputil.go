package httputil

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"diplom/internal/apperror"
	"diplom/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

const maxBodySize = 1 << 20 // 1 MiB

func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		WriteError(w, http.StatusBadRequest, response.CodeBadRequest, "invalid JSON body")
		return false
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		WriteError(w, http.StatusBadRequest, response.CodeBadRequest, "request body must be a single JSON object")
		return false
	}
	return true
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, code, message string) {
	WriteJSON(w, status, response.Error(code, message))
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	var ve *apperror.ValidationError
	switch {
	case errors.Is(err, apperror.ErrNotFound):
		WriteError(w, http.StatusNotFound, response.CodeNotFound, "resource not found")
	case errors.As(err, &ve):
		WriteError(w, http.StatusBadRequest, response.CodeBadRequest, ve.Message)
	case errors.Is(err, apperror.ErrConflict):
		WriteError(w, http.StatusConflict, response.CodeConflict, "resource already exists")
	default:
		reqID := middleware.GetReqID(r.Context())
		WriteError(w, http.StatusInternalServerError, response.CodeInternal,
			"internal server error (request_id="+reqID+")")
	}
}

func URLParamUUID(w http.ResponseWriter, r *http.Request, name string) (uuid.UUID, bool) {
	raw := chi.URLParam(r, name)
	if raw == "" {
		WriteError(w, http.StatusBadRequest, response.CodeBadRequest, "missing path parameter: "+name)
		return uuid.Nil, false
	}
	id, err := uuid.Parse(raw)
	if err != nil {
		WriteError(w, http.StatusBadRequest, response.CodeBadRequest, "invalid UUID: "+name)
		return uuid.Nil, false
	}
	return id, true
}

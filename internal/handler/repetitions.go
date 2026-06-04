package handler

import (
	"net/http"

	"diplom/internal/dto"
	"diplom/internal/service"
	"diplom/pkg/httputil"
)

type RepetitionHandler struct {
	svc *service.RepetitionService
}

func NewRepetitionHandler(svc *service.RepetitionService) *RepetitionHandler {
	return &RepetitionHandler{svc: svc}
}

func (h *RepetitionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := httputil.URLParamUUID(w, r, "id")
	if !ok {
		return
	}
	userID, ok := httputil.UserIDFromContext(w, r)
	if !ok {
		return
	}
	var req dto.UpdateRepetitionRequest
	if !httputil.DecodeJSON(w, r, &req) {
		return
	}
	rep, err := h.svc.UpdateStatus(id, userID, req)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, rep)
}

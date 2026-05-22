package handler

import (
	"net/http"

	"diplom/internal/dto"
	"diplom/internal/service"
	"diplom/pkg/httputil"
)

type SubjectHandler struct {
	svc *service.SubjectService
}

func NewSubjectHandler(svc *service.SubjectService) *SubjectHandler {
	return &SubjectHandler{svc: svc}
}

func (h *SubjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := httputil.URLParamUUID(w, r, "id")
	if !ok {
		return
	}
	var req dto.CreateSubjectRequest
	if !httputil.DecodeJSON(w, r, &req) {
		return
	}
	subject, err := h.svc.Create(userID, req)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusCreated, subject)
}

func (h *SubjectHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := httputil.URLParamUUID(w, r, "id")
	if !ok {
		return
	}
	list, err := h.svc.ListByUser(userID)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, list)
}

func (h *SubjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := httputil.URLParamUUID(w, r, "id")
	if !ok {
		return
	}
	subject, err := h.svc.Get(id)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, subject)
}

func (h *SubjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := httputil.URLParamUUID(w, r, "id")
	if !ok {
		return
	}
	var req dto.UpdateSubjectRequest
	if !httputil.DecodeJSON(w, r, &req) {
		return
	}
	subject, err := h.svc.Update(id, req)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, subject)
}

func (h *SubjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := httputil.URLParamUUID(w, r, "id")
	if !ok {
		return
	}
	if err := h.svc.Delete(id); err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

package handler

import (
	"net/http"

	"diplom/internal/dto"
	"diplom/internal/service"
	"diplom/pkg/httputil"
)

type TopicHandler struct {
	svc *service.TopicService
}

func NewTopicHandler(svc *service.TopicService) *TopicHandler {
	return &TopicHandler{svc: svc}
}

func (h *TopicHandler) Create(w http.ResponseWriter, r *http.Request) {
	subjectID, ok := httputil.URLParamUUID(w, r, "id")
	if !ok {
		return
	}
	var req dto.CreateTopicRequest
	if !httputil.DecodeJSON(w, r, &req) {
		return
	}
	topic, err := h.svc.Create(subjectID, req)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusCreated, topic)
}

func (h *TopicHandler) ListBySubject(w http.ResponseWriter, r *http.Request) {
	subjectID, ok := httputil.URLParamUUID(w, r, "id")
	if !ok {
		return
	}
	list, err := h.svc.ListBySubject(subjectID)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, list)
}

func (h *TopicHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := httputil.URLParamUUID(w, r, "id")
	if !ok {
		return
	}
	topic, err := h.svc.Get(id)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, topic)
}

func (h *TopicHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := httputil.URLParamUUID(w, r, "id")
	if !ok {
		return
	}
	var req dto.UpdateTopicRequest
	if !httputil.DecodeJSON(w, r, &req) {
		return
	}
	topic, err := h.svc.Update(id, req)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, topic)
}

func (h *TopicHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

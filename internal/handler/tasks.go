package handler

import (
	"net/http"

	"diplom/internal/dto"
	"diplom/internal/service"
	"diplom/pkg/httputil"
)

type TaskHandler struct {
	svc *service.TaskService
}

func NewTaskHandler(svc *service.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	topicID, ok := httputil.URLParamUUID(w, r, "id")
	if !ok {
		return
	}
	userID, ok := httputil.UserIDFromContext(w, r)
	if !ok {
		return
	}
	var req dto.CreateTaskRequest
	if !httputil.DecodeJSON(w, r, &req) {
		return
	}
	task, err := h.svc.Create(topicID, userID, req)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) ListByTopic(w http.ResponseWriter, r *http.Request) {
	topicID, ok := httputil.URLParamUUID(w, r, "id")
	if !ok {
		return
	}
	userID, ok := httputil.UserIDFromContext(w, r)
	if !ok {
		return
	}
	list, err := h.svc.ListByTopic(topicID, userID)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, list)
}

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := httputil.URLParamUUID(w, r, "id")
	if !ok {
		return
	}
	userID, ok := httputil.UserIDFromContext(w, r)
	if !ok {
		return
	}
	task, err := h.svc.Get(id, userID)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := httputil.URLParamUUID(w, r, "id")
	if !ok {
		return
	}
	userID, ok := httputil.UserIDFromContext(w, r)
	if !ok {
		return
	}
	var req dto.UpdateTaskRequest
	if !httputil.DecodeJSON(w, r, &req) {
		return
	}
	task, err := h.svc.Update(id, userID, req)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := httputil.URLParamUUID(w, r, "id")
	if !ok {
		return
	}
	userID, ok := httputil.UserIDFromContext(w, r)
	if !ok {
		return
	}
	if err := h.svc.Delete(id, userID); err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

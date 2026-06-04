package handler

import (
	"net/http"

	"diplom/internal/dto"
	"diplom/internal/service"
	"diplom/pkg/httputil"
)

type SessionHandler struct {
	svc *service.SessionService
}

func NewSessionHandler(svc *service.SessionService) *SessionHandler {
	return &SessionHandler{svc: svc}
}

func (h *SessionHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := httputil.UserIDFromContext(w, r)
	if !ok {
		return
	}
	var req dto.CreateSessionRequest
	if !httputil.DecodeJSON(w, r, &req) {
		return
	}
	session, err := h.svc.Create(userID, req)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusCreated, session)
}

func (h *SessionHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := httputil.URLParamUUID(w, r, "id")
	if !ok {
		return
	}
	userID, ok := httputil.UserIDFromContext(w, r)
	if !ok {
		return
	}
	session, err := h.svc.Get(id, userID)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, session)
}

func (h *SessionHandler) Finish(w http.ResponseWriter, r *http.Request) {
	id, ok := httputil.URLParamUUID(w, r, "id")
	if !ok {
		return
	}
	userID, ok := httputil.UserIDFromContext(w, r)
	if !ok {
		return
	}
	var req dto.FinishSessionRequest
	if r.ContentLength > 0 {
		if !httputil.DecodeJSON(w, r, &req) {
			return
		}
	}
	stats, err := h.svc.Finish(id, userID, req)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, stats)
}

func (h *SessionHandler) SubmitAttempt(w http.ResponseWriter, r *http.Request) {
	id, ok := httputil.URLParamUUID(w, r, "id")
	if !ok {
		return
	}
	userID, ok := httputil.UserIDFromContext(w, r)
	if !ok {
		return
	}
	var req dto.CreateTaskAttemptRequest
	if !httputil.DecodeJSON(w, r, &req) {
		return
	}
	resp, err := h.svc.SubmitAttempt(id, userID, req)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusCreated, resp)
}

func (h *SessionHandler) ListAttempts(w http.ResponseWriter, r *http.Request) {
	id, ok := httputil.URLParamUUID(w, r, "id")
	if !ok {
		return
	}
	userID, ok := httputil.UserIDFromContext(w, r)
	if !ok {
		return
	}
	list, err := h.svc.ListAttempts(id, userID)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, list)
}

func (h *SessionHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	id, ok := httputil.URLParamUUID(w, r, "id")
	if !ok {
		return
	}
	userID, ok := httputil.UserIDFromContext(w, r)
	if !ok {
		return
	}
	tasks, err := h.svc.GetTasksForSession(id, userID)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, tasks)
}

func (h *SessionHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	id, ok := httputil.URLParamUUID(w, r, "id")
	if !ok {
		return
	}
	userID, ok := httputil.UserIDFromContext(w, r)
	if !ok {
		return
	}
	stats, err := h.svc.GetStats(id, userID)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, stats)
}

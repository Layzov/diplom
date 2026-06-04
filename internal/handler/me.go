package handler

import (
	"net/http"

	"diplom/internal/apperror"
	"diplom/internal/auth"
	"diplom/internal/dto"
	"diplom/internal/models"
	"diplom/internal/service"
	"diplom/pkg/httputil"
	"diplom/pkg/response"
)

type MeHandler struct {
	users       *service.UserService
	subjects    *service.SubjectService
	sessions    *service.SessionService
	repetitions *service.RepetitionService
	statistics  *service.StatisticsService
}

func NewMeHandler(svc *service.Services) *MeHandler {
	return &MeHandler{
		users:       svc.Users,
		subjects:    svc.Subjects,
		sessions:    svc.Sessions,
		repetitions: svc.Repetitions,
		statistics:  svc.Statistics,
	}
}

func (h *MeHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := httputil.UserIDFromContext(w, r)
	if !ok {
		return
	}
	user, err := h.users.Get(userID)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, user)
}

func (h *MeHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := httputil.UserIDFromContext(w, r)
	if !ok {
		return
	}
	var req dto.UpdateUserRequest
	if !httputil.DecodeJSON(w, r, &req) {
		return
	}
	user, err := h.users.Update(userID, req)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, user)
}

func (h *MeHandler) ListSubjects(w http.ResponseWriter, r *http.Request) {
	userID, ok := httputil.UserIDFromContext(w, r)
	if !ok {
		return
	}
	list, err := h.subjects.ListByUser(userID)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, list)
}

func (h *MeHandler) CreateSubject(w http.ResponseWriter, r *http.Request) {
	userID, ok := httputil.UserIDFromContext(w, r)
	if !ok {
		return
	}
	var req dto.CreateSubjectRequest
	if !httputil.DecodeJSON(w, r, &req) {
		return
	}
	subject, err := h.subjects.Create(userID, req)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusCreated, subject)
}

func (h *MeHandler) ListSessions(w http.ResponseWriter, r *http.Request) {
	userID, ok := httputil.UserIDFromContext(w, r)
	if !ok {
		return
	}
	list, err := h.sessions.ListByUser(userID)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, list)
}

func (h *MeHandler) ListRepetitions(w http.ResponseWriter, r *http.Request) {
	userID, ok := httputil.UserIDFromContext(w, r)
	if !ok {
		return
	}
	var status *models.RepetitionStatus
	if raw := r.URL.Query().Get("status"); raw != "" {
		s := models.RepetitionStatus(raw)
		status = &s
	}
	list, err := h.repetitions.ListByUser(userID, status)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, list)
}

func (h *MeHandler) Calendar(w http.ResponseWriter, r *http.Request) {
	userID, ok := httputil.UserIDFromContext(w, r)
	if !ok {
		return
	}
	from, to, ok := parseTimeRange(w, r)
	if !ok {
		return
	}
	list, err := h.repetitions.Calendar(userID, from, to)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, list)
}

func (h *MeHandler) StatsOverview(w http.ResponseWriter, r *http.Request) {
	userID, ok := httputil.UserIDFromContext(w, r)
	if !ok {
		return
	}
	stats, err := h.statistics.UserOverview(userID)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, stats)
}

func (h *MeHandler) StatsTopics(w http.ResponseWriter, r *http.Request) {
	userID, ok := httputil.UserIDFromContext(w, r)
	if !ok {
		return
	}
	list, err := h.statistics.TopicsByUser(userID)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, list)
}

func (h *MeHandler) StatsSessions(w http.ResponseWriter, r *http.Request) {
	userID, ok := httputil.UserIDFromContext(w, r)
	if !ok {
		return
	}
	list, err := h.statistics.SessionsByUser(userID)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, list)
}

func (h *MeHandler) StatsUpcoming(w http.ResponseWriter, r *http.Request) {
	userID, ok := httputil.UserIDFromContext(w, r)
	if !ok {
		return
	}
	limit, includeOverdue, ok := parseUpcomingQuery(w, r)
	if !ok {
		return
	}
	list, err := h.statistics.UpcomingRepetitions(userID, limit, includeOverdue)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, list)
}

func requireAdmin(w http.ResponseWriter, r *http.Request) bool {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		httputil.HandleError(w, r, apperror.ErrUnauthorized)
		return false
	}
	if claims.Role != models.UserRoleAdmin {
		httputil.WriteError(w, http.StatusForbidden, response.CodeForbidden, "admin role required")
		return false
	}
	return true
}

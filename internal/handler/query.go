package handler

import (
	"net/http"
	"strconv"
	"time"

	"diplom/internal/service"
	"diplom/pkg/httputil"
	"diplom/pkg/response"
)

func parseTimeRange(w http.ResponseWriter, r *http.Request) (from, to *time.Time, ok bool) {
	from, err := parseQueryTime(r, "from")
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, response.CodeBadRequest, "invalid from: use RFC3339")
		return nil, nil, false
	}
	to, err = parseQueryTime(r, "to")
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, response.CodeBadRequest, "invalid to: use RFC3339")
		return nil, nil, false
	}
	return from, to, true
}

func parseQueryTime(r *http.Request, key string) (*time.Time, error) {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return nil, err
	}
	utc := t.UTC()
	return &utc, nil
}

func parseUpcomingQuery(w http.ResponseWriter, r *http.Request) (limit int, includeOverdue bool, ok bool) {
	limit, err := service.ParseUpcomingLimit(r.URL.Query().Get("limit"))
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return 0, false, false
	}
	includeOverdue, err = parseBoolQuery(r.URL.Query().Get("include_overdue"), false)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, response.CodeBadRequest, "include_overdue must be true or false")
		return 0, false, false
	}
	return limit, includeOverdue, true
}

func parseBoolQuery(raw string, defaultVal bool) (bool, error) {
	if raw == "" {
		return defaultVal, nil
	}
	return strconv.ParseBool(raw)
}

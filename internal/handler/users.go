package handler

import (
	"net/http"

	"diplom/internal/service"
	"diplom/pkg/httputil"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// List is admin-only.
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	list, err := h.svc.List()
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, list)
}

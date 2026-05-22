package handler

import (
	"net/http"

	"diplom/internal/dto"
	"diplom/internal/service"
	"diplom/pkg/httputil"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if !httputil.DecodeJSON(w, r, &req) {
		return
	}
	user, err := h.svc.Create(req)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.List()
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, list)
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := httputil.URLParamUUID(w, r, "id")
	if !ok {
		return
	}
	user, err := h.svc.Get(id)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := httputil.URLParamUUID(w, r, "id")
	if !ok {
		return
	}
	var req dto.UpdateUserRequest
	if !httputil.DecodeJSON(w, r, &req) {
		return
	}
	user, err := h.svc.Update(id, req)
	if err != nil {
		httputil.HandleError(w, r, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

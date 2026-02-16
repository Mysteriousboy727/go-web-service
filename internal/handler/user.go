package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"go-industry-server/internal/model"
	"go-industry-server/internal/repository"
	"go-industry-server/internal/service"
	"go-industry-server/pkg/response"
)

type UserHandler struct {
	service service.UserService
	logger  *slog.Logger
}

func NewUserHandler(svc service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{service: svc, logger: logger}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.service.CreateUser(r.Context(), req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusCreated, user)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	user, err := h.service.GetUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "user not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}
	response.Success(w, http.StatusOK, user)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list users")
		return
	}
	response.Success(w, http.StatusOK, users)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req model.UpdateUserRequest
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.service.UpdateUser(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "user not found")
			return
		}
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusOK, user)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.service.DeleteUser(r.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "user not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}
	response.Success(w, http.StatusOK, map[string]string{"message": "deleted"})
}

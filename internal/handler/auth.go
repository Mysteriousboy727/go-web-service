package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"go-industry-server/internal/auth"
	"go-industry-server/internal/repository"
	"go-industry-server/pkg/response"
)

type AuthHandler struct {
	userRepo repository.UserRepository
	jwtSvc   *auth.JWTService
	logger   *slog.Logger
}

func NewAuthHandler(userRepo repository.UserRepository, jwtSvc *auth.JWTService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
		jwtSvc:   jwtSvc,
		logger:   logger,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	user, err := h.userRepo.GetByEmail(r.Context(), req.Email)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		response.Error(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	accessToken, refreshToken, err := h.jwtSvc.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "token generation failed")
		return
	}

	response.Success(w, http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}


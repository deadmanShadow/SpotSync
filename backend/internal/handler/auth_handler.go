package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"spotsync-backend/internal/dto"
	"spotsync-backend/internal/service"
	"spotsync-backend/pkg/utils"
)

// AuthHandler exposes HTTP handlers for the authentication endpoints.
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler wires an AuthHandler bound to the given AuthService.
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles POST /api/v1/auth/register.
func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	if err := c.Validate(&req); err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	user, err := h.authService.Register(req)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			return utils.JSONError(c, http.StatusBadRequest, "Email already exists", err.Error())
		}
		return utils.JSONError(c, http.StatusInternalServerError, "Failed to register user", err.Error())
	}

	return utils.JSONSuccess(c, http.StatusCreated, "User registered successfully", user)
}

// Login handles POST /api/v1/auth/login.
func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	if err := c.Validate(&req); err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	resp, err := h.authService.Login(req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return utils.JSONError(c, http.StatusUnauthorized, "Invalid email or password", err.Error())
		}
		return utils.JSONError(c, http.StatusInternalServerError, "Failed to login", err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, "Login successful", resp)
}
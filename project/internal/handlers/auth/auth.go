package auth

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"project/internal/config"
	"project/internal/logger"
	"project/internal/models"
	"project/internal/services"
	"project/internal/services/auth"
)

type Handler struct {
	services  *services.Services
	validator *validator.Validate
	config    *config.Config
}

func NewHandler(s *services.Services, v *validator.Validate, cfg *config.Config) *Handler {
	return &Handler{
		services:  s,
		validator: v,
		config:    cfg,
	}
}

type AuthResponse struct {
	AccessToken string       `json:"access_token"`
	User        *models.User `json:"user"`
}
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// Login
// @Summary Логін
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/signin [post]
func (h *Handler) Login(c fiber.Ctx) error {
	var req LoginRequest

	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "invalid request data",
		})
	}
	err := h.validator.Struct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "invalid credentials",
		})
	}

	res, err := h.services.Auth.Login(c.Context(), auth.LoginInput{
		Email:     req.Email,
		Password:  req.Password,
		UserAgent: c.Get("User-Agent"),
		IP:        c.IP(),
	})

	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) || errors.Is(err, models.ErrInvalidPassword) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "invalid credentials",
			})
		}
		logger.LogWarnContext(c.Context(), err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "internal server error",
		})
	}

	ex := time.Now().Add(h.config.RefreshExpiration)

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		HTTPOnly: true,
		Secure:   h.config.IsProduction,
		SameSite: "Lax",
		Expires:  ex,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "device_session_id",
		Value:    res.DeviceSessionID,
		HTTPOnly: true,
		Secure:   h.config.IsProduction,
		SameSite: "Lax",
		Expires:  ex,
	})
	return c.Status(fiber.StatusOK).JSON(
		AuthResponse{
			AccessToken: res.AccessToken,
			User:        res.User,
		},
	)
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// Register
// @Summary Реєстрація
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /api/v1/auth/signup [post]
func (h *Handler) Register(c fiber.Ctx) error {
	var req RegisterRequest

	err := c.Bind().JSON(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request data",
		})
	}

	err = h.validator.Struct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	user, err := h.services.Auth.Register(c.Context(), auth.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, models.ErrUserDuplicate) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "user already exists",
			})
		}
		logger.LogWarnContext(c.Context(), err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(
		AuthResponse{
			User: user,
		},
	)
}

// Refresh
// @Summary Оновити токен
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/refresh [post]
func (h *Handler) Refresh(c fiber.Ctx) error {
	token := c.Cookies("refresh_token")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "refresh token not found",
		})
	}
	deviceSessionID := c.Cookies("device_session_id")
	if deviceSessionID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing device session",
		})
	}
	output, err := h.services.Auth.Refresh(c.Context(), auth.RefreshInput{
		RefreshToken:    token,
		DeviceSessionID: deviceSessionID,
		UserAgent:       c.Get("User-Agent"),
		IP:              c.IP(),
	})
	if err != nil {
		if errors.Is(err, models.ErrSessionNotFound) ||
			errors.Is(err, models.ErrSessionExpired) ||
			errors.Is(err, models.ErrInvalidDevice) ||
			errors.Is(err, models.ErrSuspiciousActivity) {
			c.ClearCookie("refresh_token")
			c.ClearCookie("device_session_id")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "unauthorized",
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "internal server error",
		})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    output.RefreshToken,
		HTTPOnly: true,
		Secure:   h.config.IsProduction,
		SameSite: "Lax",
		Expires:  time.Now().Add(h.config.RefreshExpiration),
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": output.AccessToken,
	})
}

// Logout
// @Summary Вийти з поточного девайсу
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/logout [post]
func (h *Handler) Logout(c fiber.Ctx) error {
	deviceSessionID := c.Cookies("device_session_id")
	if deviceSessionID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing device session",
		})
	}

	if err := h.services.Auth.Logout(c.Context(), deviceSessionID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	c.ClearCookie("refresh_token")
	c.ClearCookie("device_session_id")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "logged out successfully",
	})
}

// LogoutAll
// @Summary Вийти з усіх девайсів
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/logout/all [post]
func (h *Handler) LogoutAll(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	if err := h.services.Auth.LogoutAll(c.Context(), userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	c.ClearCookie("refresh_token")
	c.ClearCookie("device_session_id")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "logged out from all devices",
	})
}

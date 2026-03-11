package user

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"project/internal/models"
	"project/internal/services"
)

type Handler struct {
	service   *services.Services
	validator *validator.Validate
}

func NewHandler(s *services.Services, v *validator.Validate) *Handler {
	return &Handler{
		service:   s,
		validator: v,
	}
}

func (h *Handler) GetUsers(c fiber.Ctx) error {
	return nil
}

func (h *Handler) GetUser(c fiber.Ctx) error {
	return nil
}

// UsersRequest represents batch users request
// @Description List of user IDs for batch fetch
type UsersRequest struct {
	Ids []string `json:"ids"`
}

// GetUsersByIds godoc
// @Summary      Get users by IDs
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      UsersRequest            true  "List of user IDs"
// @Success      200      {object}  map[string]models.User
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /users/batch [post]
// @Security     BearerAuth
func (h *Handler) GetUsersByIds(c fiber.Ctx) error {
	var req UsersRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	users, err := h.service.User.GetByIds(c.Context(), req.Ids)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"users": users,
	})
}

// Search godoc
// @Summary      Search users
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        search  query     string          true  "Search query (username or email)"
// @Success      200     {object}  map[string][]models.User
// @Failure      400     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /users/search [get]
// @Security     BearerAuth
func (h *Handler) Search(c fiber.Ctx) error {
	query := c.Query("search")
	userID := c.Locals("user_id").(string)
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "search query is required",
		})
	}

	users, err := h.service.User.Search(c.Context(), userID, query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"users": users,
	})
}

func (h *Handler) GetProfile(c fiber.Ctx) error {
	userId := c.Locals("user_id").(string)

	usr, err := h.service.User.FindByID(c.Context(), userId)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": usr,
	})
}

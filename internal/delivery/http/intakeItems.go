package http
import (
	"strconv"
	"github.com/gofiber/fiber/v3"
	"pvz-service/internal/service"
)
type IntakeItemHandler struct {
	service *service.IntakeItemService
}
func NewIntakeItemHandler(s *service.IntakeItemService) *IntakeItemHandler {
	return &IntakeItemHandler{service: s}
}
func (h *IntakeItemHandler) AddItem(c fiber.Ctx) error {
	pvzID, err := strconv.Atoi(c.Params("id"))
	if err != nil || pvzID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid PVZ ID"})
	}
	var req struct {
		ProductID int     `json:"product_id"`
		Quantity  int     `json:"quantity"`
		Price     float64 `json:"price"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	item, err := h.service.AddItem(c.Context(), pvzID, req.ProductID, req.Quantity, req.Price)
	if err != nil {
		if err == service.ErrNoActiveIntake {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "no active intake"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(item)
}
func (h *IntakeItemHandler) DeleteLastItem(c fiber.Ctx) error {
	pvzID, err := strconv.Atoi(c.Params("id"))
	if err != nil || pvzID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid PVZ ID"})
	}
	err = h.service.DeleteLastItem(c.Context(), pvzID)
	if err != nil {
		if err == service.ErrNoActiveIntake {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "no active intake"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
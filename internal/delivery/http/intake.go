package http
import (
	"strconv"
	"github.com/gofiber/fiber/v3"
	"pvz-service/internal/service"
)
type IntakeHandler struct {
	service *service.IntakeService
}
func NewIntakeHandler(s *service.IntakeService) *IntakeHandler {
	return &IntakeHandler{service: s}
}
func (h *IntakeHandler) CreateIntake(c fiber.Ctx) error {
	type request struct {
		PVZID int `json:"pvz_id"`
	}
	var req request
	if err := c.Bind().Body(&req); err != nil || req.PVZID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request",
		})
	}
	intake, err := h.service.CreateIntake(c.Context(), req.PVZID)
	if err != nil {
		if err == service.ErrIntakeAlreadyOpen {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "intake already in progress for this PVZ",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create intake",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(intake)
}
func (h *IntakeHandler) CloseIntake(c fiber.Ctx) error {
	pvzID, err := strconv.Atoi(c.Params("id"))
	if err != nil || pvzID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid PVZ ID"})
	}
	intake, err := h.service.CloseIntake(c.Context(), pvzID)
	if err != nil {
		if err == service.ErrNoActiveIntake {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "no active intake"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(intake)
}
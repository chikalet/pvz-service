package http
import (
	"github.com/gofiber/fiber/v3"
	"pvz-service/internal/service"
	"strconv"
	"time"
)
type PVZHandler struct {
	service *service.PVZService
}
func NewPVZHandler(s *service.PVZService) *PVZHandler {
	return &PVZHandler{service: s}
}
func (h *PVZHandler) CreatePVZ(c fiber.Ctx) error {
	type request struct {
		City string `json:"city"`
	}
	var req request
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request",
			"details": err.Error(),
		})
	}
	if req.City == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "City is required",
		})
	}
	pvz, err := h.service.CreatePVZ(c.Context(), req.City)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(pvz)
}
func (h *PVZHandler) GetPVZ(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}
	pvz, err := h.service.GetPVZ(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "PVZ not found",
		})
	}
	return c.JSON(pvz)
}
func (h *PVZHandler) GetPVZs(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	var startDate, endDate *time.Time
	if startStr := c.Query("start_date"); startStr != "" {
		if t, err := time.Parse(time.RFC3339, startStr); err == nil {
			startDate = &t
		}
	}
	if endStr := c.Query("end_date"); endStr != "" {
		if t, err := time.Parse(time.RFC3339, endStr); err == nil {
			endDate = &t
		}
	}
	pvzs, err := h.service.GetPVZs(c.Context(), service.PVZFilter{
		StartDate: startDate,
		EndDate:   endDate,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(pvzs)
}
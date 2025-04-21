package http
import (
	"github.com/gofiber/fiber/v3"
	"pvz-service/internal/auth"
)
func RegisterAuthRoutes(router fiber.Router) {
	router.Get("/dummyLogin", func(c fiber.Ctx) error {
		role := c.Query("role")
		if role == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "role is required",
			})
		}
		token, err := auth.GenerateToken(role)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to generate token",
			})
		}
		return c.JSON(fiber.Map{
			"token": token,
			"role":  role,
		})
	})
}
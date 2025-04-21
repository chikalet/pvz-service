package auth
import (
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"pvz-service/metrics"
	"time"
)
func RequireRole(requiredRole string) fiber.Handler {
	return func(c fiber.Ctx) error {
		tokenStr := c.Get("Authorization")
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Требуется авторизация",
			})
		}
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte("supersecret"), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Неверный токен",
			})
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Ошибка чтения прав доступа",
			})
		}
		role, ok := claims["role"].(string)
		if !ok || role != requiredRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Доступ запрещён для вашей роли",
			})
		}
		return c.Next()
	}
}
func AuthRequired(c fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Токен отсутствует",
		})
	}
	_, err := ParseToken(token)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Невалидный токен",
		})
	}
	return c.Next()
}
func MetricsMiddleware(c fiber.Ctx) error {
	start := time.Now()
	path := c.Path()
	method := c.Method()
	err := c.Next()
	status := c.Response().StatusCode()
	metrics.HTTPRequests.WithLabelValues(path, method, string(status)).Inc()
	metrics.HTTPDuration.WithLabelValues(path).Observe(time.Since(start).Seconds())
	return err
}
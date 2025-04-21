package main
import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"pvz-service/internal/auth"
	authMiddleware "pvz-service/internal/auth"
	"pvz-service/internal/config"
	httpDelivery "pvz-service/internal/delivery/http"
	"pvz-service/internal/service"
	"pvz-service/internal/storage/postgres"
)
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	cfg := config.Load()
	log.Println("Config loaded")
	db, err := postgres.NewStorage(ctx, cfg.DB.ConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	app := fiber.New(fiber.Config{
		AppName:      "PVZ Service",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${ip} ${status} - ${method} ${path}\n",
	}))
	pvzService := service.NewPVZService(db)
	intakeService := service.NewIntakeService(db)
	itemService := service.NewIntakeItemService(db, db)
	itemHandler := httpDelivery.NewIntakeItemHandler(itemService)
	pvzHandler := httpDelivery.NewPVZHandler(pvzService)
	intakeHandler := httpDelivery.NewIntakeHandler(intakeService)
	api := app.Group("/api/v1")
	api.Post("/pvz/:id/items", authMiddleware.RequireRole("employee"), itemHandler.AddItem)
	api.Delete("/pvz/:id/items", authMiddleware.RequireRole("employee"), itemHandler.DeleteLastItem)
	httpDelivery.RegisterAuthRoutes(api)
	api.Post("/pvz", authMiddleware.RequireRole("moderator"), pvzHandler.CreatePVZ)
	api.Get("/pvz/:id", pvzHandler.GetPVZ)
	api.Post("/intake", authMiddleware.RequireRole("employee"), intakeHandler.CreateIntake)
	api.Post("/pvz/:id/close", authMiddleware.RequireRole("employee"), intakeHandler.CloseIntake)
	api.Get("/pvz", authMiddleware.RequireRole("moderator"), pvzHandler.GetPVZs)
	app.Use(func(c fiber.Ctx) error {
		if c.Path() == "/api/v1/dummyLogin" {
			return c.Next()
		}
		return authMiddleware.AuthRequired(c)
	})
	api.Get("/pvz", authMiddleware.RequireRole("moderator"), pvzHandler.GetPVZs)
	api.Post("/pvz", authMiddleware.RequireRole("moderator"), pvzHandler.CreatePVZ)
	go func() {
		log.Printf("Server is running on port %s\n", cfg.Server.Port)
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()
	startMetricsServer()
	app.Use(auth.MetricsMiddleware)
	<-ctx.Done()
	log.Println("Shutting down...")
	if err := app.ShutdownWithTimeout(5 * time.Second); err != nil {
		log.Fatalf("Shutdown failed: %v", err)
	}
	log.Println("Server stopped cleanly.")
}
func startMetricsServer() {
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		srv := &http.Server{
			Addr:    ":9000",
			Handler: mux,
		}
		log.Println("Prometheus metrics server started on :9000")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("metrics server error: %v", err)
		}
	}()
}
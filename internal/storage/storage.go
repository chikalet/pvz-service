package storage
import (
	"context"
	"errors"
	"pvz-service/internal/models"
	"time"
)
var ErrNotFound = errors.New("not found")
type PVZRepository interface {
	CreatePVZ(ctx context.Context, city string) (*models.PVZ, error)
	GetPVZByID(ctx context.Context, id string) (*models.PVZ, error)
	GetPVZs(ctx context.Context, filter PVZFilter) ([]*models.PVZ, error)
}
type IntakeRepository interface {
	GetOpenIntake(ctx context.Context, pvzID int) (*models.Intake, error)
	CreateIntake(ctx context.Context, pvzID int) (*models.Intake, error)
	CloseIntake(ctx context.Context, pvzID int) (*models.Intake, error)
}
type IntakeItemRepository interface {
	AddItemToIntake(ctx context.Context, intakeID, productID, quantity int, price float64) (*models.IntakeItem, error)
	DeleteLastItem(ctx context.Context, intakeID int) error
}
type PVZFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
	Limit     int
	Offset    int
}
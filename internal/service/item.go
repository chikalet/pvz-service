package service
import (
	"context"
	"errors"
	"pvz-service/internal/models"
	"pvz-service/internal/storage"
	"pvz-service/metrics"
)
type IntakeItemService struct {
	itemRepo   storage.IntakeItemRepository
	intakeRepo storage.IntakeRepository
}
var ErrNoActiveIntake = errors.New("no active intake for this PVZ")
func NewIntakeItemService(itemRepo storage.IntakeItemRepository, intakeRepo storage.IntakeRepository) *IntakeItemService {
	return &IntakeItemService{itemRepo: itemRepo, intakeRepo: intakeRepo}
}
func (s *IntakeItemService) AddItem(ctx context.Context, pvzID, productID, quantity int, price float64) (*models.IntakeItem, error) {
	intake, err := s.intakeRepo.GetOpenIntake(ctx, pvzID)
	if err != nil {
		if err == storage.ErrNotFound {
			return nil, ErrNoActiveIntake
		}
		return nil, err
	}
	metrics.ItemsAdded.Inc()                     
	metrics.ItemsQuantity.Add(float64(quantity)) 
	return s.itemRepo.AddItemToIntake(ctx, intake.ID, productID, quantity, price)
}
func (s *IntakeItemService) DeleteLastItem(ctx context.Context, pvzID int) error {
	intake, err := s.intakeRepo.GetOpenIntake(ctx, pvzID)
	if err != nil {
		if err == storage.ErrNotFound {
			return ErrNoActiveIntake
		}
		return err
	}
	return s.itemRepo.DeleteLastItem(ctx, intake.ID)
}
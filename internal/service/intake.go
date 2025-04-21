package service
import (
	"context"
	"errors"
	"pvz-service/internal/models"
	"pvz-service/internal/storage"
	"pvz-service/metrics"
)
var ErrIntakeAlreadyOpen = errors.New("intake already in progress")
type IntakeService struct {
	repo storage.IntakeRepository
}
func NewIntakeService(r storage.IntakeRepository) *IntakeService {
	return &IntakeService{repo: r}
}
func (s *IntakeService) CreateIntake(ctx context.Context, pvzID int) (*models.Intake, error) {
	_, err := s.repo.GetOpenIntake(ctx, pvzID)
	if err == nil {
		return nil, ErrIntakeAlreadyOpen
	}
	if err != storage.ErrNotFound {
		return nil, err
	}
	metrics.IntakesOpened.Inc()
	return s.repo.CreateIntake(ctx, pvzID)
}
func (s *IntakeService) CloseIntake(ctx context.Context, pvzID int) (*models.Intake, error) {
	_, err := s.repo.GetOpenIntake(ctx, pvzID)
	if err != nil {
		if err == storage.ErrNotFound {
			return nil, ErrNoActiveIntake
		}
		return nil, err
	}
	metrics.IntakesClosed.Inc()
	return s.repo.CloseIntake(ctx, pvzID)
}
package service
import (
	"context"
	"errors"
	"pvz-service/internal/models"
	"pvz-service/internal/storage"
	"pvz-service/metrics" 
	"time"
)
var (
	ErrInvalidCity = errors.New("only Moscow, SPb and Kazan are allowed")
)
type PVZFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
	Limit     int
	Offset    int
}
type PVZService struct {
	repo storage.PVZRepository
}
func NewPVZService(repo storage.PVZRepository) *PVZService {
	return &PVZService{repo: repo}
}
func (s *PVZService) CreatePVZ(ctx context.Context, city string) (*models.PVZ, error) {
	allowedCities := map[string]struct{}{
		"Москва":          {},
		"Санкт-Петербург": {},
		"Казань":          {},
	}
	if _, ok := allowedCities[city]; !ok {
		return nil, ErrInvalidCity
	}
	pvz, err := s.repo.CreatePVZ(ctx, city)
	if err != nil {
		return nil, err
	}
	metrics.PVZCreated.Inc() 
	return pvz, nil
}
func (s *PVZService) GetPVZ(ctx context.Context, id string) (*models.PVZ, error) {
	return s.repo.GetPVZByID(ctx, id)
}
func (s *PVZService) GetPVZs(ctx context.Context, filter PVZFilter) ([]*models.PVZ, error) {
	if filter.Limit <= 0 || filter.Limit > 100 {
		filter.Limit = 10
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	return s.repo.GetPVZs(ctx, storage.PVZFilter{
		StartDate: filter.StartDate,
		EndDate:   filter.EndDate,
		Limit:     filter.Limit,
		Offset:    filter.Offset,
	})
}
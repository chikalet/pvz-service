package postgres
import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"pvz-service/internal/models"
	"pvz-service/internal/storage"
)
func (s *Storage) CreatePVZ(ctx context.Context, city string) (*models.PVZ, error) {
	const query = `
	INSERT INTO pickup_points (city) 
	VALUES ($1)
	RETURNING id, created_at`
	var pvz models.PVZ
	err := s.db.QueryRow(ctx, query, city).Scan(&pvz.ID, &pvz.CreatedAt)
	if err != nil {
		return nil, err
	}
	pvz.City = city
	return &pvz, nil
}
func (s *Storage) GetPVZByID(ctx context.Context, id string) (*models.PVZ, error) {
	const query = `
	SELECT id, city, created_at 
	FROM pickup_points 
	WHERE id = $1`
	var pvz models.PVZ
	err := s.db.QueryRow(ctx, query, id).Scan(&pvz.ID, &pvz.City, &pvz.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, storage.ErrNotFound
	}
	return &pvz, err
}
func (s *Storage) GetPVZs(ctx context.Context, filter storage.PVZFilter) ([]*models.PVZ, error) {
	query := `
        SELECT DISTINCT p.id, p.city, p.created_at
        FROM pickup_points p
        LEFT JOIN intakes i ON p.id = i.pvz_id
        WHERE 1=1
    `
	args := []interface{}{}
	argPos := 1
	if filter.StartDate != nil {
		query += fmt.Sprintf(" AND i.started_at >= $%d", argPos)
		args = append(args, *filter.StartDate)
		argPos++
	}
	if filter.EndDate != nil {
		query += fmt.Sprintf(" AND i.started_at <= $%d", argPos)
		args = append(args, *filter.EndDate)
		argPos++
	}
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, filter.Limit, filter.Offset)
	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pvzs []*models.PVZ
	for rows.Next() {
		var pvz models.PVZ
		if err := rows.Scan(&pvz.ID, &pvz.City, &pvz.CreatedAt); err != nil {
			return nil, err
		}
		pvzs = append(pvzs, &pvz)
	}
	return pvzs, nil
}
package postgres
import (
	"context"
	"github.com/jackc/pgx/v5"
	"pvz-service/internal/models"
	"pvz-service/internal/storage"
)
func (s *Storage) GetOpenIntake(ctx context.Context, pvzID int) (*models.Intake, error) {
	const query = `
		SELECT id, pvz_id, status, started_at, closed_at
		FROM intakes
		WHERE pvz_id = $1 AND status = 'in_progress'
		LIMIT 1`
	var intake models.Intake
	err := s.db.QueryRow(ctx, query, pvzID).Scan(
		&intake.ID,
		&intake.PVZID,
		&intake.Status,
		&intake.StartedAt,
		&intake.ClosedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, storage.ErrNotFound
	}
	return &intake, err
}
func (s *Storage) CreateIntake(ctx context.Context, pvzID int) (*models.Intake, error) {
	const query = `
		INSERT INTO intakes (pvz_id, status)
		VALUES ($1, 'in_progress')
		RETURNING id, started_at`
	var intake models.Intake
	intake.PVZID = pvzID
	intake.Status = models.IntakeInProgress
	err := s.db.QueryRow(ctx, query, pvzID).Scan(&intake.ID, &intake.StartedAt)
	return &intake, err
}
func (s *Storage) CloseIntake(ctx context.Context, pvzID int) (*models.Intake, error) {
	const query = `
        UPDATE intakes
        SET status = 'closed',
            closed_at = NOW()
        WHERE pvz_id = $1 AND status = 'in_progress'
        RETURNING id, pvz_id, status, started_at, closed_at`
	var intake models.Intake
	err := s.db.QueryRow(ctx, query, pvzID).Scan(
		&intake.ID,
		&intake.PVZID,
		&intake.Status,
		&intake.StartedAt,
		&intake.ClosedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, storage.ErrNotFound
	}
	return &intake, err
}
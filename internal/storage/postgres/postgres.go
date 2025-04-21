package postgres
import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)
type Storage struct {
	db *pgxpool.Pool
}
func NewStorage(ctx context.Context, dsn string) (*Storage, error) {
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}
func (s *Storage) Close() {
	s.db.Close()
}
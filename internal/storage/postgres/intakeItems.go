package postgres
import (
	"context"
	"time"
	"pvz-service/internal/models"
)
func (s *Storage) AddItemToIntake(ctx context.Context, intakeID, productID, quantity int, price float64) (*models.IntakeItem, error) {
	const query = `
		INSERT INTO intake_items (intake_id, product_id, quantity, price, received_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, received_at`
	var item models.IntakeItem
	item.IntakeID = intakeID
	item.ProductID = productID
	item.Quantity = quantity
	item.Price = price
	now := time.Now()
	err := s.db.QueryRow(ctx, query, intakeID, productID, quantity, price, now).Scan(&item.ID, &item.ReceivedAt)
	return &item, err
}
func (s *Storage) DeleteLastItem(ctx context.Context, intakeID int) error {
	const query = `
		DELETE FROM intake_items
		WHERE id = (
			SELECT id FROM intake_items
			WHERE intake_id = $1
			ORDER BY received_at DESC
			LIMIT 1
		);`
	_, err := s.db.Exec(ctx, query, intakeID)
	return err
}
package models
import (
	"time"
)
type IntakeItem struct {
	ID         int       `json:"id"`
	IntakeID   int       `json:"intake_id"`
	ProductID  int       `json:"product_id"`
	Quantity   int       `json:"quantity"`
	Price      float64   `json:"price"`
	ReceivedAt time.Time `json:"received_at"`
}
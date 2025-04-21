package models
import "time"
type PVZ struct {
	ID        int       `json:"id"`
	City      string    `json:"city"`
	CreatedAt time.Time `json:"created_at"`
}
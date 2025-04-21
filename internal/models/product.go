package models
import "time"
type ProductType string
const (
	ProductElectronics ProductType = "electronics"
	ProductClothing    ProductType = "clothing"
	ProductShoes       ProductType = "shoes"
)
type Product struct {
	ID       int         `json:"id"`
	IntakeID int         `json:"intake_id"`
	Type     ProductType `json:"type"`
	AddedAt  time.Time   `json:"added_at"`
}
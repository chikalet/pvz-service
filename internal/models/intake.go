package models
import "time"
type IntakeStatus string
const (
	IntakeInProgress IntakeStatus = "in_progress"
	IntakeClosed     IntakeStatus = "closed"
)
type Intake struct {
	ID        int          `json:"id"`
	PVZID     int          `json:"pvz_id"`
	Status    IntakeStatus `json:"status"`
	StartedAt time.Time    `json:"started_at"`
	ClosedAt  *time.Time   `json:"closed_at,omitempty"`
}
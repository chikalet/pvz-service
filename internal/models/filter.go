package models
import "time"
type PVZFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
	Limit     int
	Offset    int
}
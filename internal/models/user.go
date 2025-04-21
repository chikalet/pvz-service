package models
import "time"
type UserRole string
const (
	RoleClient    UserRole = "client"
	RoleModerator UserRole = "moderator"
	RoleEmployee  UserRole = "employee"
)
type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}
package model

// User represents a user in the system
type User struct {
	BaseModel
	Username     string `json:"username" gorm:"uniqueIndex;not null"`
	Email        string `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string `json:"-" gorm:"not null"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Role         string `json:"role"`
	IsActive     bool   `json:"is_active" gorm:"default:true"`
}

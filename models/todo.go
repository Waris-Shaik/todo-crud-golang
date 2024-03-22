package models

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Title       string   `json:"title" gorm:"not null"`
	Description string   `json:"description"`
	Completed   bool     `json:"completed" gorm:"default:false"`
	UserID      uint     `json:"user_id"`                                 // Foreign Key for the user model
	User        UserLite `json:"user,omitempty" gorm:"foreignKey:UserID"` // User association
}
type UserLite struct {
	ID       uint   `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
}

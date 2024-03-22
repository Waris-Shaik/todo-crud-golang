package models

import "time"

type User struct {
	ID        uint      `json:"_id" gorm:"primarykey"`
	Name      string    `json:"name"`
	UserName  string    `json:"username"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:null"`
}

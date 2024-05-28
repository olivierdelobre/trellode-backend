package models

import "time"

type User struct {
	ID           int       `gorm:"column:id;primaryKey" json:"id"`
	Email        string    `gorm:"column:email" json:"email"`
	Password     string    `gorm:"-" json:"password"`
	PasswordHash string    `gorm:"column:password_hash" json:"-"`
	CreatedAt    time.Time `gorm:"created_at" json:"createdAt"`
}

func (User) TableName() string {
	return "users"
}
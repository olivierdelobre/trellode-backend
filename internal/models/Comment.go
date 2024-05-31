package models

import "time"

type Comment struct {
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	CardID    string    `gorm:"column:card_id" json:"cardId"`
	UserID    string    `gorm:"column:user_id" json:"userId"`
	Content   string    `gorm:"column:content" json:"content"`
	CreatedAt time.Time `gorm:"created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"updated_at" json:"updatedAt"`
}

func (Comment) TableName() string {
	return "comments"
}

package models

import "time"

type Background struct {
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	UserID    string    `gorm:"column:user_id" json:"userId"`
	Data      string    `gorm:"column:data" json:"data"`
	Color     string    `gorm:"color" json:"color"`
	CreatedAt time.Time `gorm:"created_at" json:"createdAt"`
}

func (Background) TableName() string {
	return "backgrounds"
}

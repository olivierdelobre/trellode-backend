package models

import "time"

type Background struct {
	ID     int `gorm:"column:id;primaryKey" json:"id"`
	UserID int `gorm:"column:user_id" json:"userId"`
	//Data       []byte    `gorm:"column:data" json:"-"`
	Data       string    `gorm:"column:data" json:"data"`
	DataBase64 string    `gorm:"-" json:"dataBase64"`
	Color      string    `gorm:"color" json:"color"`
	CreatedAt  time.Time `gorm:"created_at" json:"createdAt"`
}

func (Background) TableName() string {
	return "backgrounds"
}

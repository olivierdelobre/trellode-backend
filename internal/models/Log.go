package models

import "time"

type Log struct {
	ID                string    `gorm:"column:id;primaryKey" json:"id"`
	UserID            string    `gorm:"column:user_id" json:"userId"`
	User              *User     `gorm:"foreignKey:UserID" json:"user"`
	BoardID           string    `gorm:"column:board_id" json:"boardId"`
	Action            string    `gorm:"column:action" json:"action"`
	ActionTargetID    string    `gorm:"column:action_target_id" json:"actionTargetId"`
	ActionTargetTitle string    `gorm:"-" json:"actionTargetTitle"`
	Changes           string    `gorm:"column:changes" json:"changes"` // json structure containing what has changed
	CreatedAt         time.Time `gorm:"created_at" json:"createdAt"`
}

func (Log) TableName() string {
	return "logs"
}

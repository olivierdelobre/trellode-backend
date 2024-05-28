package models

import "time"

type Board struct {
	ID              int        `gorm:"column:id;primaryKey" json:"id"`
	UserID          int        `gorm:"column:user_id" json:"userId"`
	Title           string     `gorm:"column:title" json:"title"`
	BackgroundImage string     `gorm:"column:background_image" json:"backgroundImage"`
	Lists           []List     `gorm:"foreignKey:BoardID" json:"lists"`
	CreatedAt       time.Time  `gorm:"created_at" json:"createdAt"`
	UpdatedAt       time.Time  `gorm:"updated_at" json:"updatedAt"`
	ArchivedAt      *time.Time `gorm:"archived_at" json:"archivedAt"`
}

func (Board) TableName() string {
	return "boards"
}

package models

import "time"

type Board struct {
	ID             string      `gorm:"column:id;primaryKey" json:"id"`
	UserID         string      `gorm:"column:user_id" json:"userId"`
	Title          string      `gorm:"column:title" json:"title"`
	BackgroundID   string      `gorm:"column:background_id" json:"backgroundId"`
	Background     *Background `gorm:"foreignKey:BackgroundID" json:"background"`
	MenuColorLight string      `gorm:"-" json:"menuColorLight"`
	MenuColorDark  string      `gorm:"-" json:"menuColorDark"`
	ListColor      string      `gorm:"-" json:"listColor"`
	Lists          []List      `gorm:"foreignKey:BoardID" json:"lists"`
	CreatedAt      time.Time   `gorm:"created_at" json:"createdAt"`
	UpdatedAt      time.Time   `gorm:"updated_at" json:"updatedAt"`
	ArchivedAt     *time.Time  `gorm:"archived_at" json:"archivedAt"`
	OpenedAt       time.Time   `gorm:"opened_at" json:"openedAt"`
}

func (Board) TableName() string {
	return "boards"
}

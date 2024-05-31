package models

import "time"

type Checklist struct {
	ID         string          `gorm:"column:id;primaryKey" json:"id"`
	CardID     string          `gorm:"column:card_id" json:"cardId"`
	Title      string          `gorm:"column:title" json:"title"`
	Items      []ChecklistItem `gorm:"foreignKey:ChecklistID" json:"items"`
	CreatedAt  time.Time       `gorm:"created_at" json:"createdAt"`
	UpdatedAt  time.Time       `gorm:"updated_at" json:"updatedAt"`
	ArchivedAt *time.Time      `gorm:"archived_at" json:"archivedAt"`
}

func (Checklist) TableName() string {
	return "checklists"
}

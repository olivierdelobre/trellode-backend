package models

import "time"

type ChecklistItem struct {
	ID          string    `gorm:"column:id;primaryKey" json:"id"`
	Title       string    `gorm:"column:title" json:"title"`
	ChecklistID string    `gorm:"column:checklist_id" json:"checklistId"`
	Position    int       `gorm:"column:position" json:"position"`
	Checked     bool      `gorm:"column:checked" json:"checked"`
	CreatedAt   time.Time `gorm:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"updated_at" json:"updatedAt"`
}

func (ChecklistItem) TableName() string {
	return "checklistitems"
}

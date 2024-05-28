package models

import "time"

type Card struct {
	ID          int        `gorm:"column:id;primaryKey" json:"id"`
	ListID      int        `gorm:"column:list_id" json:"listId"`
	Title       string     `gorm:"column:title" json:"title"`
	Description string     `gorm:"column:description" json:"description"`
	Position    int        `gorm:"column:position" json:"position"`
	Comments    []Comment  `gorm:"foreignKey:CardID" json:"comments"`
	CreatedAt   time.Time  `gorm:"created_at" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"updated_at" json:"updatedAt"`
	ArchivedAt  *time.Time `gorm:"archived_at" json:"archivedAt"`
}

func (Card) TableName() string {
	return "cards"
}

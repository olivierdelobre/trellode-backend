package models

import "time"

type List struct {
	ID         string     `gorm:"column:id;primaryKey" json:"id"`
	BoardID    string     `gorm:"column:board_id" json:"boardId"`
	Title      string     `gorm:"column:title" json:"title"`
	Position   int        `gormjson:"position"`
	Cards      []Card     ` gorm:"foreignKey:ListID" json:"cards"`
	CreatedAt  time.Time  `gorm:"created_at" json:"createdAt"`
	UpdatedAt  time.Time  `gorm:"updated_at" json:"updatedAt"`
	ArchivedAt *time.Time `gorm:"archived_at" json:"archivedAt"`
}

func (balise *List) TableName() string {
	return "lists"
}

package models

import (
	"time"

	"gorm.io/gorm"
)

// Vertex reprezentuje wierzchołek grafu (mikroserwis)
type Vertex struct {
	ID          string         `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description,omitempty"`
	ParentID    *string        `json:"parent_id,omitempty" gorm:"index"` // ID rodzica (null = root)
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName określa nazwę tabeli w bazie danych
func (Vertex) TableName() string {
	return "vertices"
}

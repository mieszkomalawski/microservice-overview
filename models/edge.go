package models

import (
	"time"

	"gorm.io/gorm"
)

// Edge reprezentuje relację między wierzchołkami
type Edge struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	From      string    `json:"from" gorm:"not null;index"`
	To        string    `json:"to" gorm:"not null;index"`
	Type      string    `json:"type,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName określa nazwę tabeli w bazie danych
func (Edge) TableName() string {
	return "edges"
}



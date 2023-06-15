package models

import (
	"gorm.io/gorm"
	"time"
)

type ProjectCore struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	AuthorID  uint
	User      UserCore `gorm:""`
	Json      string   `gorm:"not null;size: 65535"`
}

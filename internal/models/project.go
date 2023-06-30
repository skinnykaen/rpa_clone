package models

import (
	"gorm.io/gorm"
	"time"
)

type ProjectCore struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	AuthorID  uint
	User      UserCore `gorm:"foreignKey:AuthorID"`
	IsShared  bool
	Json      string `gorm:"not null;size: 65535" json:"json"`
}

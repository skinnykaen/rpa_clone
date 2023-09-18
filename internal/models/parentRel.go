package models

import (
	"gorm.io/gorm"
	"time"
)

// TODO set gorm table name (ParentCore)

type ParentRelCore struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ParentID  uint
	Parent    UserCore `gorm:"foreignKey:ParentID"`
	ChildID   uint
	Child     UserCore `gorm:"foreignKey:ChildID"`
}

package models

import (
	"gorm.io/gorm"
	"time"
)

type RobboGroupRelCore struct {
	ID           uint `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	UserID       uint
	User         UserCore `gorm:"foreignKey:UserID"`
	RobboGroupID uint
	RobboGroup   RobboGroupCore `gorm:"foreignKey:RobboGroupID"`
}

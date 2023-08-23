package models

import (
	"gorm.io/gorm"
	"time"
)

type RobboGroupCore struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string         `gorm:"size:256;not null"`
	RobboUnitID uint
	RobboUnit   RobboUnitCore `gorm:"foreignKey:RobboUnitID"`
}

package models

import (
	"gorm.io/gorm"
	"time"
)

type RobboUnitRelCore struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	UnitAdminID uint
	UnitAdmin   UserCore `gorm:"foreignKey:UnitAdminID"`
	RobboUnitID uint
	RobboUnit   RobboUnitCore `gorm:"foreignKey:RobboUnitID"`
}

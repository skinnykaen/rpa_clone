package models

import (
	"gorm.io/gorm"
	"time"
)

type ProjectPageCore struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	//FIXME foreignKey
	ProjectID     uint
	Project       ProjectCore
	Title         string `gorm:"size:256;not null"`
	Instruction   string `gorm:"size:256;not null"`
	Notes         string `gorm:"size:256;not null"`
	LinkToScratch string `gorm:"size:256;not null"`
	IsShared      bool
}

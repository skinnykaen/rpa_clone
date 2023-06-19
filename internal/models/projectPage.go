package models

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

type ProjectPageCore struct {
	ID            uint `gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	AuthorID      uint
	User          UserCore `gorm:"foreignKey:AuthorID"`
	ProjectID     uint
	Project       ProjectCore `gorm:"foreignKey:ProjectID"`
	Title         string      `gorm:"size:256;not null"`
	Instruction   string      `gorm:"size:256;not null"`
	Notes         string      `gorm:"size:256;not null"`
	LinkToScratch string      `gorm:"size:256;not null"`
	IsShared      bool
}

func (p *ProjectPageHTTP) FromCore(projectPage ProjectPageCore) {
	p.ID = strconv.Itoa(int(projectPage.ID))
	p.CreatedAt = projectPage.CreatedAt.String()
	p.UpdatedAt = projectPage.UpdatedAt.String()
	p.AuthorID = strconv.Itoa(int(projectPage.AuthorID))
	p.ProjectID = strconv.Itoa(int(projectPage.ProjectID))
	p.ProjectUpdatedAt = projectPage.Project.UpdatedAt.String()
	p.Title = projectPage.Title
	p.Instruction = projectPage.Instruction
	p.Notes = projectPage.Notes
	p.LinkToScratch = projectPage.LinkToScratch
	p.IsShared = projectPage.IsShared
}

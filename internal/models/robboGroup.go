package models

import (
	"gorm.io/gorm"
	"strconv"
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

func (r *RobboGroupHTTP) FromCore(robboGroup RobboGroupCore) {
	r.ID = strconv.Itoa(int(robboGroup.ID))
	r.CreatedAt = robboGroup.CreatedAt.Format(time.DateTime)
	r.UpdatedAt = robboGroup.UpdatedAt.Format(time.DateTime)
	r.Name = robboGroup.Name
	r.RobboUnit = &RobboUnitHTTP{}
	r.RobboUnit.FromCore(robboGroup.RobboUnit)
}

func FromRobboGroupsCore(robboGroupsCore []RobboGroupCore) (robboGroupsHttp []*RobboGroupHTTP) {
	for _, robboGroupCore := range robboGroupsCore {
		var tmpRobboGroupHttp RobboGroupHTTP
		tmpRobboGroupHttp.FromCore(robboGroupCore)
		robboGroupsHttp = append(robboGroupsHttp, &tmpRobboGroupHttp)
	}
	return robboGroupsHttp
}

package models

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

type RobboUnitCore struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"size:256;not null"`
	City      string         `gorm:"size:256;not null"`
}

func (r *RobboUnitHTTP) FromCore(robboUnit RobboUnitCore) {
	r.ID = strconv.Itoa(int(robboUnit.ID))
	r.CreatedAt = robboUnit.CreatedAt.Format(time.DateTime)
	r.UpdatedAt = robboUnit.UpdatedAt.Format(time.DateTime)
	r.Name = robboUnit.Name
	r.City = robboUnit.City
}

func FromRobboUnitsCore(robboUnitsCore []RobboUnitCore) (robboUnitsHttp []*RobboUnitHTTP) {
	for _, robboUnitCore := range robboUnitsCore {
		var tmpRobboUnitHttp RobboUnitHTTP
		tmpRobboUnitHttp.FromCore(robboUnitCore)
		robboUnitsHttp = append(robboUnitsHttp, &tmpRobboUnitHttp)
	}
	return robboUnitsHttp
}

package models

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

type UserCore struct {
	ID             uint `gorm:"primaryKey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	Email          string         `gorm:"not null;"`
	Password       string         `gorm:"not null;"`
	Role           Role           `gorm:"not null;"`
	Firstname      string         `gorm:"not null;"`
	Lastname       string         `gorm:"not null;"`
	Middlename     string         `gorm:"not null;"`
	Nickname       string         `gorm:"not null;"`
	IsActive       bool           `gorm:"not null;default:false"`
	ActivationCode uint           `gorm:"not null;"`
}

func (u *UserHTTP) ToCore() UserCore {
	id, _ := strconv.ParseUint(u.ID, 10, 64)
	return UserCore{
		ID:             uint(id),
		Email:          u.Email,
		Password:       u.Password,
		Role:           u.Role,
		Firstname:      u.Firstname,
		Lastname:       u.Lastname,
		Middlename:     u.Middlename,
		Nickname:       u.Nickname,
		IsActive:       u.IsActive,
		ActivationCode: uint(u.ActivationCode),
	}
}

func (u *UserHTTP) FromCore(userCore UserCore) {
	u.ID = strconv.Itoa(int(userCore.ID))
	u.Email = userCore.Email
	u.Firstname = userCore.Firstname
	u.Lastname = userCore.Lastname
	u.Middlename = userCore.Middlename
	u.Nickname = userCore.Nickname
	u.IsActive = userCore.IsActive
	u.Role = userCore.Role
}

func FromUsersCore(usersCore []UserCore) (usersHttp []*UserHTTP) {
	for _, userCore := range usersCore {
		var tmpUserHttp UserHTTP
		tmpUserHttp.FromCore(userCore)
		usersHttp = append(usersHttp, &tmpUserHttp)
	}
	return
}

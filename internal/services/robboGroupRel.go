package services

import (
	"github.com/skinnykaen/rpa_clone/internal/gateways"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
)

type RobboGroupRelService interface {
	CreateRel(userId, robboGroupId uint) (models.RobboGroupRelCore, error)
	DeleteRel(userId, robboGroupId uint) (err error)
	GetRelById(id uint) (models.RobboGroupRelCore, error)
	GetStudentsByRobboGroupId(page, pageSize *int, robboGroupId uint) (students []models.UserCore, countRows int, err error)
	GetTeachersByRobboGroupId(page, pageSize *int, robboGroupId uint) (teachers []models.UserCore, countRows int, err error)
	GetRobboGroupsByUserId(userId uint) (robboGroups []models.RobboGroupCore, err error)
}

type RobboGroupRelServiceImpl struct {
	robboGroupRelGateway gateways.RobboGroupRelGateway
}

func (r RobboGroupRelServiceImpl) CreateRel(userId, robboGroupId uint) (models.RobboGroupRelCore, error) {
	return r.robboGroupRelGateway.CreateRel(models.RobboGroupRelCore{UserID: userId, RobboGroupID: robboGroupId})
}

func (r RobboGroupRelServiceImpl) DeleteRel(userId, robboGroupId uint) (err error) {
	return r.robboGroupRelGateway.DeleteRel(models.RobboGroupRelCore{UserID: userId, RobboGroupID: robboGroupId})
}

func (r RobboGroupRelServiceImpl) GetRelById(id uint) (models.RobboGroupRelCore, error) {
	return r.robboGroupRelGateway.GetRelById(id)
}

func (r RobboGroupRelServiceImpl) GetStudentsByRobboGroupId(page, pageSize *int, robboGroupId uint) (students []models.UserCore, countRows int, err error) {
	offset, limit := utils.GetOffsetAndLimit(page, pageSize)
	return r.robboGroupRelGateway.GetStudentsByRobboGroupId(offset, limit, robboGroupId)
}

func (r RobboGroupRelServiceImpl) GetTeachersByRobboGroupId(page, pageSize *int, robboGroupId uint) (teachers []models.UserCore, countRows int, err error) {
	offset, limit := utils.GetOffsetAndLimit(page, pageSize)
	return r.robboGroupRelGateway.GetTeachersByRobboGroupId(offset, limit, robboGroupId)
}

func (r RobboGroupRelServiceImpl) GetRobboGroupsByUserId(userId uint) (robboGroups []models.RobboGroupCore, err error) {
	return r.robboGroupRelGateway.GetRobboGroupsByUserId(userId)
}

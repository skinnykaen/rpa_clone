package services

import (
	"github.com/skinnykaen/rpa_clone/internal/gateways"
	"github.com/skinnykaen/rpa_clone/internal/models"
)

type RobboUnitRelService interface {
	CreateRel(unitAdminId, robboUnitId uint) (models.RobboUnitRelCore, error)
	DeleteRel(unitAdminId, robboUnitId uint) (err error)
	GetRelById(id uint) (models.RobboUnitRelCore, error)
	GetUnitAdminsByRobboUnitId(robboUnitId uint) (unitAdmins []models.UserCore, err error)
	GetRobboUnitsByUnitAdmin(unitAdminId uint) (robboUnits []models.RobboUnitCore, err error)
	GetStudentsByRobboUnitId(robboUnitId uint) (students []models.UserCore, err error)
}

type robboGroupsByUnitAdminProvider interface {
	GetRobboGroupsByRobboUnitById(offset, limit int, robboUnitId uint) (robboGroups []models.RobboGroupCore, countRows uint, err error)
}

type RobboUnitRelServiceImpl struct {
	robboUnitRelGateway            gateways.RobboUnitRelGateway
	robboGroupsByUnitAdminProvider robboGroupsByUnitAdminProvider
}

func (u RobboUnitRelServiceImpl) GetStudentsByRobboUnitId(robboUnitId uint) (students []models.UserCore, err error) {
	return u.robboUnitRelGateway.GetStudentsByRobboUnitId(robboUnitId)
}

func (u RobboUnitRelServiceImpl) CreateRel(unitAdminId, robboUnitId uint) (models.RobboUnitRelCore, error) {
	return u.robboUnitRelGateway.CreateRel(models.RobboUnitRelCore{UnitAdminID: unitAdminId, RobboUnitID: robboUnitId})
}

func (u RobboUnitRelServiceImpl) DeleteRel(unitAdminId, robboUnitId uint) (err error) {
	return u.robboUnitRelGateway.DeleteRel(models.RobboUnitRelCore{UnitAdminID: unitAdminId, RobboUnitID: robboUnitId})
}

func (u RobboUnitRelServiceImpl) GetRelById(id uint) (models.RobboUnitRelCore, error) {
	return u.robboUnitRelGateway.GetRelById(id)
}

func (u RobboUnitRelServiceImpl) GetUnitAdminsByRobboUnitId(robboUnitId uint) (unitAdmins []models.UserCore, err error) {
	return u.robboUnitRelGateway.GetUnitAdminsByRobboUnitId(robboUnitId)
}

func (u RobboUnitRelServiceImpl) GetRobboUnitsByUnitAdmin(unitAdminId uint) (robboUnits []models.RobboUnitCore, err error) {
	return u.robboUnitRelGateway.GetRobboUnitsByUnitAdmin(unitAdminId)
}

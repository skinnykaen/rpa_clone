package services

import (
	"github.com/skinnykaen/rpa_clone/internal/gateways"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
)

type RobboUnitService interface {
	CreateRobboUnit(robboUnit models.RobboUnitCore) (models.RobboUnitCore, error)
	GetRobboUnitById(id uint) (robboUnit models.RobboUnitCore, err error)
	DeleteRobboUnit(id uint) error
	UpdateRobboUnit(robboUnit models.RobboUnitCore) (updated models.RobboUnitCore, err error)
	GetAllRobboUnits(page, pageSize *int, clientRole models.Role) (robboUnits []models.RobboUnitCore, countRows uint, err error)
}

type RobboUnitServiceImpl struct {
	robboUnitGateway gateways.RobboUnitGateway
}

func (r RobboUnitServiceImpl) DeleteRobboUnit(id uint) error {
	return r.robboUnitGateway.DeleteRobboUnit(id)
}

func (r RobboUnitServiceImpl) UpdateRobboUnit(robboUnit models.RobboUnitCore) (updated models.RobboUnitCore, err error) {
	return r.robboUnitGateway.UpdateRobboUnit(robboUnit)
}

func (r RobboUnitServiceImpl) GetAllRobboUnits(page, pageSize *int, clientRole models.Role) (robboUnits []models.RobboUnitCore, countRows uint, err error) {
	offset, limit := utils.GetOffsetAndLimit(page, pageSize)
	if clientRole.String() == models.RoleSuperAdmin.String() {
		return r.robboUnitGateway.GetAllRobboUnits(offset, limit)
	}
	// TODO robbo units for unit admin
	return
}

func (r RobboUnitServiceImpl) GetRobboUnitById(id uint) (robboUnit models.RobboUnitCore, err error) {
	return r.robboUnitGateway.GetRobboUnitById(id)
}

func (r RobboUnitServiceImpl) CreateRobboUnit(robboUnit models.RobboUnitCore) (models.RobboUnitCore, error) {
	return r.robboUnitGateway.CreateRobboUnit(robboUnit)
}

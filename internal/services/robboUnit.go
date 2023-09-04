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
	GetAllRobboUnits(page, pageSize *int, clientId uint, clientRole models.Role) (robboUnits []models.RobboUnitCore, countRows uint, err error)
}

type RobboUnitServiceImpl struct {
	robboUnitGateway              gateways.RobboUnitGateway
	robboUnitsByUnitAdminProvider robboUnitsByUnitAdminProvider
}

func (r RobboUnitServiceImpl) DeleteRobboUnit(id uint) error {
	return r.robboUnitGateway.DeleteRobboUnit(id)
}

func (r RobboUnitServiceImpl) UpdateRobboUnit(robboUnit models.RobboUnitCore) (updated models.RobboUnitCore, err error) {
	return r.robboUnitGateway.UpdateRobboUnit(robboUnit)
}

func (r RobboUnitServiceImpl) GetAllRobboUnits(page, pageSize *int, clientId uint, clientRole models.Role) (robboUnits []models.RobboUnitCore, countRows uint, err error) {
	offset, limit := utils.GetOffsetAndLimit(page, pageSize)
	switch clientRole.String() {
	case models.RoleSuperAdmin.String():
		return r.robboUnitGateway.GetAllRobboUnits(offset, limit)
	case models.RoleUnitAdmin.String():
		robboUnits, err := r.robboUnitsByUnitAdminProvider.GetRobboUnitsByUnitAdmin(clientId)
		if err != nil {
			return []models.RobboUnitCore{}, 0, err
		}
		return robboUnits, uint(len(robboUnits)), nil
	}
	return
}

func (r RobboUnitServiceImpl) GetRobboUnitById(id uint) (robboUnit models.RobboUnitCore, err error) {
	return r.robboUnitGateway.GetRobboUnitById(id)
}

func (r RobboUnitServiceImpl) CreateRobboUnit(robboUnit models.RobboUnitCore) (models.RobboUnitCore, error) {
	return r.robboUnitGateway.CreateRobboUnit(robboUnit)
}

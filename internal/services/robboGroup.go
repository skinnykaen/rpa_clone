package services

import (
	"github.com/skinnykaen/rpa_clone/internal/gateways"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
)

type RobboGroupService interface {
	CreateRobboGroup(robboGroup models.RobboGroupCore) (models.RobboGroupCore, error)
	GetRobboGroupById(id uint) (robboGroup models.RobboGroupCore, err error)
	DeleteRobboGroup(id uint) error
	UpdateRobboGroup(robboGroup models.RobboGroupCore) (models.RobboGroupCore, error)
	GetAllRobboGroups(page, pageSize *int, clientId uint, clientRole models.Role) (robboGroups []models.RobboGroupCore, countRows uint, err error)
	GetRobboGroupsByRobboUnitById(page, pageSize *int, robboUnitId uint) (robboGroups []models.RobboGroupCore, countRows uint, err error)
}

type robboGroupsByTeacherProvider interface {
	GetRobboGroupsByUserId(userId uint) (robboGroups []models.RobboGroupCore, err error)
}

type RobboGroupServiceImpl struct {
	robboGroupGateway             gateways.RobboGroupGateway
	robboUnitsByUnitAdminProvider robboUnitsByUnitAdminProvider
	robboGroupsByTeacherProvider  robboGroupsByTeacherProvider
}

func (r RobboGroupServiceImpl) GetRobboGroupsByRobboUnitById(page, pageSize *int, robboUnitId uint) (robboGroups []models.RobboGroupCore, countRows uint, err error) {
	offset, limit := utils.GetOffsetAndLimit(page, pageSize)
	return r.robboGroupGateway.GetRobboGroupsByRobboUnitById(offset, limit, robboUnitId)
}

func (r RobboGroupServiceImpl) CreateRobboGroup(robboGroup models.RobboGroupCore) (models.RobboGroupCore, error) {
	return r.robboGroupGateway.CreateRobboGroup(robboGroup)
}

func (r RobboGroupServiceImpl) GetRobboGroupById(id uint) (robboGroup models.RobboGroupCore, err error) {
	return r.robboGroupGateway.GetRobboGroupById(id)
}

func (r RobboGroupServiceImpl) DeleteRobboGroup(id uint) error {
	return r.robboGroupGateway.DeleteRobboGroup(id)
}

func (r RobboGroupServiceImpl) UpdateRobboGroup(robboGroup models.RobboGroupCore) (models.RobboGroupCore, error) {
	return r.robboGroupGateway.UpdateRobboGroup(robboGroup)
}

func (r RobboGroupServiceImpl) GetAllRobboGroups(page, pageSize *int, clientId uint, clientRole models.Role) (robboGroups []models.RobboGroupCore, countRows uint, err error) {
	offset, limit := utils.GetOffsetAndLimit(page, pageSize)
	switch clientRole.String() {
	case models.RoleSuperAdmin.String():
		return r.robboGroupGateway.GetAllRobboGroups(offset, limit)
	case models.RoleUnitAdmin.String():
		robboUnits, err := r.robboUnitsByUnitAdminProvider.GetRobboUnitsByUnitAdmin(clientId)
		if err != nil {
			return []models.RobboGroupCore{}, 0, err
		}
		for _, robboUnit := range robboUnits {
			robboGroupsByRobboUnit, _, err := r.robboGroupGateway.GetRobboGroupsByRobboUnitById(0, 100, robboUnit.ID)
			if err != nil {
				return []models.RobboGroupCore{}, 0, err
			}
			robboGroups = append(robboGroups, robboGroupsByRobboUnit...)
		}
		return robboGroups, uint(len(robboGroups)), nil
	case models.RoleTeacher.String():
		robboGroups, err := r.robboGroupsByTeacherProvider.GetRobboGroupsByUserId(clientId)
		if err != nil {
			return []models.RobboGroupCore{}, 0, err
		}
		return robboGroups, uint(len(robboGroups)), nil
	}

	return robboGroups, countRows, nil
}

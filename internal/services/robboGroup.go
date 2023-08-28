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
	GetAllRobboGroups(page, pageSize *int, clientRole models.Role) (robboGroups []models.RobboGroupCore, countRows uint, err error)
}

type RobboGroupServiceImpl struct {
	robboGroupGateway gateways.RobboGroupGateway
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

func (r RobboGroupServiceImpl) GetAllRobboGroups(page, pageSize *int, clientRole models.Role) (robboGroups []models.RobboGroupCore, countRows uint, err error) {
	offset, limit := utils.GetOffsetAndLimit(page, pageSize)
	// TODO for unit admin get groups in own unit
	return r.robboGroupGateway.GetAllRobboGroups(offset, limit)
}

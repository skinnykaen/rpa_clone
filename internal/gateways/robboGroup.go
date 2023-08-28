package gateways

import (
	"github.com/skinnykaen/rpa_clone/internal/db"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
)

type RobboGroupGateway interface {
	CreateRobboGroup(robboGroup models.RobboGroupCore) (newRobboGroup models.RobboGroupCore, err error)
	GetRobboGroupById(id uint) (robboGroup models.RobboGroupCore, err error)
	DeleteRobboGroup(id uint) error
	UpdateRobboGroup(robboGroup models.RobboGroupCore) (models.RobboGroupCore, error)
	GetAllRobboGroups(offset, limit int) (robboGroups []models.RobboGroupCore, countRows uint, err error)
}

type RobboGroupGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (r RobboGroupGatewayImpl) CreateRobboGroup(robboGroup models.RobboGroupCore) (newRobboGroup models.RobboGroupCore, err error) {
	if err := r.postgresClient.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&robboGroup).Clauses(clause.Returning{}).Error; err != nil {
			return err
		}
		var robboUnit models.RobboUnitCore
		if err := tx.Find(&robboUnit, robboGroup.RobboUnitID).Error; err != nil {
			return err
		}
		robboGroup.RobboUnit = robboUnit
		return nil
	}); err != nil {
		return models.RobboGroupCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return robboGroup, nil
}

func (r RobboGroupGatewayImpl) GetRobboGroupById(id uint) (robboGroup models.RobboGroupCore, err error) {
	if err := r.postgresClient.Db.Preload("RobboUnit").First(&robboGroup, id).Error; err != nil {
		return robboGroup, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return robboGroup, nil
}

func (r RobboGroupGatewayImpl) DeleteRobboGroup(id uint) error {
	// TODO delete students rel all rel
	if err := r.postgresClient.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.RobboGroupCore{}, id).Error; err != nil {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		return nil
	}); err != nil {
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

func (r RobboGroupGatewayImpl) UpdateRobboGroup(robboGroup models.RobboGroupCore) (models.RobboGroupCore, error) {
	if err := r.postgresClient.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&robboGroup).Updates(&robboGroup).Error; err != nil {
			return err
		}
		if err := tx.Preload("RobboUnit").First(&robboGroup, robboGroup.ID).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return models.RobboGroupCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return robboGroup, nil
}

func (r RobboGroupGatewayImpl) GetAllRobboGroups(offset, limit int) (robboGroups []models.RobboGroupCore, countRows uint, err error) {
	var count int64
	result := r.postgresClient.Db.Preload("RobboUnit").Limit(limit).Offset(offset).Find(&robboGroups)
	if result.Error != nil {
		return []models.RobboGroupCore{}, 0, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: result.Error.Error(),
		}
	}
	result.Count(&count)
	return robboGroups, uint(count), result.Error
}

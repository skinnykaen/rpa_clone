package gateways

import (
	"github.com/skinnykaen/rpa_clone/internal/db"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
)

type RobboUnitGateway interface {
	CreateRobboUnit(robboUnit models.RobboUnitCore) (newRobboUnit models.RobboUnitCore, err error)
	GetRobboUnitById(id uint) (robboUnit models.RobboUnitCore, err error)
	DeleteRobboUnit(id uint) error
	UpdateRobboUnit(robboUnit models.RobboUnitCore) (updated models.RobboUnitCore, err error)
	GetAllRobboUnits(offset, limit int) (robboUnits []models.RobboUnitCore, countRows uint, err error)
}

type RobboUnitGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (r RobboUnitGatewayImpl) GetAllRobboUnits(offset, limit int) (robboUnits []models.RobboUnitCore, countRows uint, err error) {
	var count int64
	result := r.postgresClient.Db.Limit(limit).Offset(offset).Find(&robboUnits)
	if result.Error != nil {
		return []models.RobboUnitCore{}, 0, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: result.Error.Error(),
		}
	}
	result.Count(&count)
	return robboUnits, uint(count), result.Error
}

func (r RobboUnitGatewayImpl) DeleteRobboUnit(id uint) error {
	if err := r.postgresClient.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("robbo_unit_id = ?", id).Delete(&models.RobboGroupCore{}).Error; err != nil {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		if err := tx.Delete(&models.RobboUnitCore{}, id).Error; err != nil {
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

func (r RobboUnitGatewayImpl) UpdateRobboUnit(robboUnit models.RobboUnitCore) (updated models.RobboUnitCore, err error) {
	if err = r.postgresClient.Db.Model(&robboUnit).Clauses(clause.Returning{}).Take(&models.RobboUnitCore{}, robboUnit.ID).
		Updates(
			map[string]interface{}{
				"name": robboUnit.Name,
				"city": robboUnit.City,
			},
		).Error; err != nil {
		return models.RobboUnitCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return robboUnit, nil
}

func (r RobboUnitGatewayImpl) GetRobboUnitById(id uint) (robboUnit models.RobboUnitCore, err error) {
	if err := r.postgresClient.Db.First(&robboUnit, id).Error; err != nil {
		return robboUnit, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return robboUnit, nil
}

func (r RobboUnitGatewayImpl) CreateRobboUnit(robboUnit models.RobboUnitCore) (models.RobboUnitCore, error) {
	if err := r.postgresClient.Db.Create(&robboUnit).Clauses(clause.Returning{}).Error; err != nil {
		return models.RobboUnitCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return robboUnit, nil
}

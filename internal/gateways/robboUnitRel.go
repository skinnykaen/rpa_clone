package gateways

import (
	"github.com/skinnykaen/rpa_clone/internal/db"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
)

type RobboUnitRelGateway interface {
	CreateRel(rel models.RobboUnitRelCore) (models.RobboUnitRelCore, error)
	DeleteRel(rel models.RobboUnitRelCore) (err error)
	GetRelById(id uint) (models.RobboUnitRelCore, error)
	GetUnitAdminsByRobboUnitId(robboUnitId uint) (unitAdmins []models.UserCore, err error)
	GetRobboUnitsByUnitAdmin(unitAdminId uint) (robboUnits []models.RobboUnitCore, err error)
	GetStudentsByRobboUnitId(robboUnitId uint) (users []models.UserCore, err error)
	GetTeachersByRobboUnitId(robboUnitId uint) (users []models.UserCore, err error)
}

type RobboUnitRelGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (u RobboUnitRelGatewayImpl) GetStudentsByRobboUnitId(robboUnitId uint) (students []models.UserCore, err error) {
	if err := u.postgresClient.Db.Transaction(func(tx *gorm.DB) error {
		var robboGroups []models.RobboGroupCore
		if err := tx.Preload("RobboUnit").Where("robbo_unit_id = ?", robboUnitId).Find(&robboGroups).Error; err != nil {
			return err
		}
		for _, robboGroup := range robboGroups {
			var rels []models.RobboGroupRelCore
			if err := tx.Preload("User").Where(models.RobboGroupRelCore{RobboGroupID: robboGroup.ID}).Find(&rels).Error; err != nil {
				return err
			}
			for _, rel := range rels {
				if rel.User.Role.String() == models.RoleStudent.String() {
					students = append(students, rel.User)
				}
			}
		}
		return nil
	}); err != nil {
		return []models.UserCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return students, nil
}

func (u RobboUnitRelGatewayImpl) GetTeachersByRobboUnitId(robboUnitId uint) (teachers []models.UserCore, err error) {
	if err := u.postgresClient.Db.Transaction(func(tx *gorm.DB) error {
		var robboGroups []models.RobboGroupCore
		if err := tx.Preload("RobboUnit").Where("robbo_unit_id = ?", robboUnitId).Find(&robboGroups).Error; err != nil {
			return err
		}
		userIdsMap := make(map[uint]bool)
		for _, robboGroup := range robboGroups {
			var rels []models.RobboGroupRelCore
			if err := tx.Preload("User").Where(models.RobboGroupRelCore{RobboGroupID: robboGroup.ID}).Find(&rels).Error; err != nil {
				return err
			}
			for i := 0; i < len(rels); i++ {
				if rels[i].User.Role.String() == models.RoleTeacher.String() &&
					!userIdsMap[rels[i].UserID] {
					teachers = append(teachers, rels[i].User)
				}
				userIdsMap[rels[i].UserID] = true
			}
		}
		return nil
	}); err != nil {
		return []models.UserCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return teachers, nil
}

func (u RobboUnitRelGatewayImpl) CreateRel(rel models.RobboUnitRelCore) (models.RobboUnitRelCore, error) {
	if err := u.postgresClient.Db.Transaction(func(tx *gorm.DB) error {
		unitAdmin := models.UserCore{ID: rel.UnitAdminID}
		robboUnit := models.RobboUnitCore{ID: rel.RobboUnitID}
		if err := tx.First(&unitAdmin).Error; err != nil {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		if err := tx.First(&robboUnit).Error; err != nil {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		if unitAdmin.Role != models.RoleUnitAdmin {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: "incorrect parent or child id",
			}
		}
		if err := tx.Create(&rel).Clauses(clause.Returning{}).Error; err != nil {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		return nil
	}); err != nil {
		return models.RobboUnitRelCore{}, err
	}
	return rel, nil
}

func (u RobboUnitRelGatewayImpl) DeleteRel(rel models.RobboUnitRelCore) (err error) {
	if err := u.postgresClient.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("robbo_unit_id = ? AND unit_admin_id = ?", rel.RobboUnitID, rel.UnitAdminID).Delete(&rel).Error; err != nil {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (u RobboUnitRelGatewayImpl) GetRelById(id uint) (rel models.RobboUnitRelCore, err error) {
	if err := u.postgresClient.Db.First(&rel, id).Error; err != nil {
		return models.RobboUnitRelCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return rel, nil
}

func (u RobboUnitRelGatewayImpl) GetUnitAdminsByRobboUnitId(robboUnitId uint) (unitAdmins []models.UserCore, err error) {
	var rels []models.RobboUnitRelCore
	if err := u.postgresClient.Db.Where(models.RobboUnitRelCore{RobboUnitID: robboUnitId}).Preload("UnitAdmin").Find(&rels).Error; err != nil {
		return []models.UserCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	for _, rel := range rels {
		unitAdmins = append(unitAdmins, rel.UnitAdmin)
	}
	return unitAdmins, nil
}

func (u RobboUnitRelGatewayImpl) GetRobboUnitsByUnitAdmin(unitAdminId uint) (robboUnits []models.RobboUnitCore, err error) {
	var rels []models.RobboUnitRelCore
	if err := u.postgresClient.Db.Where(models.RobboUnitRelCore{UnitAdminID: unitAdminId}).Preload("RobboUnit").Find(&rels).Error; err != nil {
		return []models.RobboUnitCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	for _, rel := range rels {
		robboUnits = append(robboUnits, rel.RobboUnit)
	}
	return robboUnits, nil
}

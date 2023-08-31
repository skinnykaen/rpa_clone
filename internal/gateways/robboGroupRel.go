package gateways

import (
	"github.com/skinnykaen/rpa_clone/internal/db"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
)

type RobboGroupRelGateway interface {
	CreateRel(rel models.RobboGroupRelCore) (models.RobboGroupRelCore, error)
	DeleteRel(rel models.RobboGroupRelCore) (err error)
	GetRelById(id uint) (models.RobboGroupRelCore, error)
	GetStudentsByRobboGroupId(offset, limit int, robboGroupId uint) (students []models.UserCore, countRows int, err error)
	GetTeachersByRobboGroupId(offset, limit int, robboGroupId uint) (teachers []models.UserCore, countRows int, err error)
	GetRobboGroupsByUserId(userId uint) (robboGroups []models.RobboGroupCore, err error)
}

type RobboGroupRelGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (r RobboGroupRelGatewayImpl) CreateRel(rel models.RobboGroupRelCore) (models.RobboGroupRelCore, error) {
	if err := r.postgresClient.Db.Transaction(func(tx *gorm.DB) error {
		user := models.UserCore{ID: rel.UserID}
		robboGroup := models.RobboGroupCore{ID: rel.RobboGroupID}
		if err := tx.First(&user).Error; err != nil {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		if err := tx.First(&robboGroup).Error; err != nil {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		if user.Role.String() != models.RoleStudent.String() && user.Role.String() != models.RoleTeacher.String() {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: "incorrect student or teacher id",
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
		return models.RobboGroupRelCore{}, err
	}
	return rel, nil
}

func (r RobboGroupRelGatewayImpl) DeleteRel(rel models.RobboGroupRelCore) (err error) {
	if err := r.postgresClient.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("robbo_group_id = ? AND user_id = ?", rel.RobboGroupID, rel.UserID).Delete(&rel).Error; err != nil {
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

func (r RobboGroupRelGatewayImpl) GetRelById(id uint) (rel models.RobboGroupRelCore, err error) {
	if err := r.postgresClient.Db.First(&rel, id).Error; err != nil {
		return models.RobboGroupRelCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return rel, nil
}

func (r RobboGroupRelGatewayImpl) GetStudentsByRobboGroupId(offset, limit int, robboGroupId uint) (students []models.UserCore, countRows int, err error) {
	var rels []models.RobboGroupRelCore
	if err := r.postgresClient.Db.Limit(limit).Offset(offset).
		Where(models.RobboGroupRelCore{RobboGroupID: robboGroupId}).Preload("User").Find(&rels).Error; err != nil {
		return []models.UserCore{}, 0, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	for _, rel := range rels {
		if rel.User.Role.String() == models.RoleStudent.String() {
			students = append(students, rel.User)
		}
	}
	return students, len(students), nil
}

func (r RobboGroupRelGatewayImpl) GetTeachersByRobboGroupId(offset, limit int, robboGroupId uint) (teachers []models.UserCore, countRows int, err error) {
	var rels []models.RobboGroupRelCore
	if err := r.postgresClient.Db.Limit(limit).Offset(offset).
		Where(models.RobboGroupRelCore{RobboGroupID: robboGroupId}).Preload("User").Find(&rels).Error; err != nil {
		return []models.UserCore{}, 0, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	for _, rel := range rels {
		if rel.User.Role.String() == models.RoleTeacher.String() {
			teachers = append(teachers, rel.User)
		}
	}
	return teachers, len(teachers), nil
}

func (r RobboGroupRelGatewayImpl) GetRobboGroupsByUserId(userId uint) (robboGroups []models.RobboGroupCore, err error) {
	var rels []models.RobboGroupRelCore
	if err := r.postgresClient.Db.Where(models.RobboGroupRelCore{UserID: userId}).Preload("RobboGroup").Find(&rels).Error; err != nil {
		return []models.RobboGroupCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	for _, rel := range rels {
		robboGroups = append(robboGroups, rel.RobboGroup)
	}
	return robboGroups, nil
}

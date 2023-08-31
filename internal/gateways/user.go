package gateways

import (
	"errors"
	"github.com/skinnykaen/rpa_clone/internal/consts"
	"github.com/skinnykaen/rpa_clone/internal/db"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
)

type UserGateway interface {
	CreateUser(user models.UserCore) (newUser models.UserCore, err error)
	DeleteUser(id uint) (err error)
	UpdateUser(user models.UserCore) (updatedUser models.UserCore, err error)
	GetUserById(id uint) (user models.UserCore, err error)
	GetUserByActivationLink(link string) (user models.UserCore, err error)
	GetUserByEmail(email string) (user models.UserCore, err error)
	GetUsersByEmail(offset, limit int, role []models.Role, email string) (users []models.UserCore, countRows uint, err error)
	GetAllUsers(offset, limit int, isActive bool, role []models.Role) (users []models.UserCore, countRows uint, err error)
	DoesExistEmail(id uint, email string) (bool, error)
	SetIsActive(id uint, isActive bool) error
	GetStudentsByRobboUnitId(offset, limit int, robboUnitId uint) (students []models.UserCore, countRows uint, err error)
}

type UserGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (u UserGatewayImpl) GetStudentsByRobboUnitId(offset, limit int, robboUnitId uint) (students []models.UserCore, countRows uint, err error) {
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
		return []models.UserCore{}, 0, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return students, uint(len(students)), nil
}

func (u UserGatewayImpl) GetUsersByEmail(offset, limit int, role []models.Role, email string) (users []models.UserCore, countRows uint, err error) {
	var count int64
	if len(role) == 0 {
		role = append(role,
			models.RoleStudent,
			models.RoleParent,
			models.RoleTeacher,
			models.RoleUnitAdmin,
		)
	}
	result := u.postgresClient.Db.Limit(limit).Offset(offset).
		Where("email LIKE ? AND (role) IN ?", email+"%", role).Find(&users)
	if result.Error != nil {
		return []models.UserCore{}, 0, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: result.Error.Error(),
		}
	}
	result.Count(&count)
	return users, uint(count), result.Error
}

func (u UserGatewayImpl) GetUserByActivationLink(link string) (user models.UserCore, err error) {
	if err = u.postgresClient.Db.Where("activation_link = ?", link).Take(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (u UserGatewayImpl) GetUserByEmail(email string) (user models.UserCore, err error) {
	if err = u.postgresClient.Db.Where("email = ?", email).Take(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, utils.ResponseError{
				Code:    http.StatusBadRequest,
				Message: consts.ErrIncorrectPasswordOrEmail,
			}
		}
		return user, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return user, nil
}

func (u UserGatewayImpl) SetIsActive(id uint, isActive bool) error {
	var updateStruct map[string]interface{}
	if isActive {
		updateStruct = map[string]interface{}{
			"is_active":       isActive,
			"activation_link": "",
		}
	} else {
		updateStruct = map[string]interface{}{
			"is_active": isActive,
		}
	}
	return u.postgresClient.Db.First(&models.UserCore{ID: id}).Updates(updateStruct).Error
}

func (u UserGatewayImpl) DoesExistEmail(id uint, email string) (bool, error) {
	result := u.postgresClient.Db.Where("id != ? AND email = ?", id, email).
		Take(&models.UserCore{})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: result.Error.Error(),
		}
	}
	return true, nil
}

func (u UserGatewayImpl) CreateUser(user models.UserCore) (newUser models.UserCore, err error) {
	result := u.postgresClient.Db.Create(&user).Clauses(clause.Returning{})
	if result.Error != nil {
		return models.UserCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: result.Error.Error(),
		}
	}
	return user, nil
}

func (u UserGatewayImpl) DeleteUser(id uint) (err error) {
	return u.postgresClient.Db.Delete(&models.UserCore{}, id).Error
}

func (u UserGatewayImpl) UpdateUser(user models.UserCore) (models.UserCore, error) {
	if err := u.postgresClient.Db.Model(&user).Clauses(clause.Returning{}).
		Take(&models.UserCore{}, user.ID).
		Updates(map[string]interface{}{
			"email":      user.Email,
			"firstname":  user.Firstname,
			"lastname":   user.Lastname,
			"middlename": user.Middlename,
			"nickname":   user.Nickname,
		}).Error; err != nil {
		return models.UserCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return user, nil
}

func (u UserGatewayImpl) GetUserById(id uint) (user models.UserCore, err error) {
	if err := u.postgresClient.Db.First(&user, id).Error; err != nil {
		return models.UserCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return user, nil
}

func (u UserGatewayImpl) GetAllUsers(
	offset, limit int,
	isActive bool,
	role []models.Role,
) (users []models.UserCore, countRows uint, err error) {
	var count int64
	if len(role) == 0 {
		role = append(role,
			models.RoleStudent,
			models.RoleParent,
			models.RoleTeacher,
			models.RoleUnitAdmin,
		)
	}
	result := u.postgresClient.Db.Limit(limit).Offset(offset).
		Where("is_active = ? AND (role) IN ?", isActive, role).Find(&users)
	if result.Error != nil {
		return []models.UserCore{}, 0, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: result.Error.Error(),
		}
	}
	result.Count(&count)
	return users, uint(count), result.Error
}

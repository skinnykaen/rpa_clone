package gateways

import (
	"errors"
	"github.com/skinnykaen/rpa_clone/internal/db"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserGateway interface {
	CreateUser(user models.UserCore) (newUser models.UserCore, err error)
	DeleteUser(id uint) (err error)
	UpdateUser(user models.UserCore) (updatedUser models.UserCore, err error)
	GetUserById(id uint) (user models.UserCore, err error)
	GetUserByEmail(email string) (user models.UserCore, err error)
	GetAllUsers(offset, limit int, isActive bool, role []models.Role) (users []models.UserCore, countRows uint, err error)
	DoesExistEmail(id uint, email string) (bool, error)
	SetIsActive(id uint, isActive bool) error
}

type UserGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (u UserGatewayImpl) GetUserByEmail(email string) (user models.UserCore, err error) {
	if err = u.postgresClient.Db.Where("email = ?", email).Take(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (u UserGatewayImpl) SetIsActive(id uint, isActive bool) error {
	return u.postgresClient.Db.First(&models.UserCore{ID: id}).Updates(
		map[string]interface{}{
			"is_active": isActive,
		}).Error
}

func (u UserGatewayImpl) DoesExistEmail(id uint, email string) (bool, error) {
	result := u.postgresClient.Db.Where("id != ? AND email = ?", id, email).
		Take(&models.UserCore{})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}

func (u UserGatewayImpl) CreateUser(user models.UserCore) (newUser models.UserCore, err error) {
	result := u.postgresClient.Db.Create(&user).Clauses(clause.Returning{})
	return user, result.Error
}

func (u UserGatewayImpl) DeleteUser(id uint) (err error) {
	return u.postgresClient.Db.Delete(&models.UserCore{}, id).Error
}

func (u UserGatewayImpl) UpdateUser(user models.UserCore) (models.UserCore, error) {
	err := u.postgresClient.Db.Model(&user).Clauses(clause.Returning{}).
		Take(&models.UserCore{}, user.ID).
		Updates(map[string]interface{}{
			"email":      user.Email,
			"firstname":  user.Firstname,
			"lastname":   user.Lastname,
			"middlename": user.Middlename,
			"nickname":   user.Nickname,
		}).Error
	return user, err
}

func (u UserGatewayImpl) GetUserById(id uint) (user models.UserCore, err error) {
	err = u.postgresClient.Db.First(&user, id).Error
	return
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
		return []models.UserCore{}, 0, result.Error
	}
	result.Count(&count)
	return users, uint(count), result.Error
}

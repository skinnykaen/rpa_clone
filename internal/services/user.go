package services

import (
	"errors"
	"rpa_clone/internal/gateways"
	"rpa_clone/internal/models"
	"rpa_clone/pkg/logger"
	"rpa_clone/pkg/utils"
)

type UserService interface {
	CreateUser()
	DeleteUser(id uint) error
	UpdateUser(user models.UserCore) (updatedUser models.UserCore, err error)
	GetUserById(id uint) (models.UserCore, error)
	GetAllUsers(page, pageSize *int, isActive bool, role []models.Role) (users []models.UserCore, countRows uint, err error)
}

type UserServiceImpl struct {
	loggers     logger.Loggers
	userGateway gateways.UserGateway
}

func (u UserServiceImpl) CreateUser() {
	//TODO implement me
	panic("implement me")
}

func (u UserServiceImpl) DeleteUser(id uint) error {
	return u.userGateway.DeleteUser(id)
}

func (u UserServiceImpl) UpdateUser(user models.UserCore) (updatedUser models.UserCore, err error) {
	exist, err := u.userGateway.DoesExistEmail(user.ID, user.Email)
	if err != nil {
		return models.UserCore{}, err
	}
	if exist {
		return models.UserCore{}, errors.New("email already in use")
	}
	return u.userGateway.UpdateUser(user)
}

func (u UserServiceImpl) GetUserById(id uint) (models.UserCore, error) {
	return u.userGateway.GetUserById(id)
}

func (u UserServiceImpl) GetAllUsers(page, pageSize *int,
	isActive bool,
	roles []models.Role,
) (users []models.UserCore, countRows uint, err error) {
	offset, limit := utils.GetOffsetAndLimit(page, pageSize)
	return u.userGateway.GetAllUsers(offset, limit, isActive, roles)
}

package services

import (
	"github.com/skinnykaen/rpa_clone/internal/consts"
	"github.com/skinnykaen/rpa_clone/internal/gateways"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
	"net/http"
)

type UserService interface {
	CreateUser(user models.UserCore, clientRole models.Role) (newUser models.UserCore, err error)
	DeleteUser(id uint) error
	UpdateUser(user models.UserCore, clientRole models.Role) (updatedUser models.UserCore, err error)
	GetUserById(id uint, clientRole models.Role) (models.UserCore, error)
	GetAllUsers(page, pageSize *int, isActive bool, roles []models.Role, clientId uint, clientRole models.Role) (users []models.UserCore, countRows uint, err error)
	GetUsersByEmail(page, pageSize *int, roles []models.Role, email string, clientRole models.Role) (users []models.UserCore, countRows uint, err error)
	SetIsActive(id uint, isActive bool) error
	GetStudentsByUnitAdminId(unitAdminId uint) (students []models.UserCore, err error)
	GetTeachersByUnitAdminId(unitAdminId uint) (students []models.UserCore, err error)
	GetClientsByUnitAdminId(unitAdminId uint) (clients []models.UserCore, err error)
	GetStudentsByTeacherId(teacherId uint) (students []models.UserCore, err error)
}

type parentByChildIdProvider interface {
	GetParentsByChildId(childId uint) (parents []models.UserCore, err error)
}

type robboGroupRelProvider interface {
	GetStudentsByRobboGroupId(offset, limit int, robboGroupId uint) (students []models.UserCore, countRows int, err error)
	GetRobboGroupsByUserId(userId uint) (robboGroups []models.RobboGroupCore, err error)
}

type robboUnitsByUnitAdminProvider interface {
	GetRobboUnitsByUnitAdmin(unitAdminId uint) (robboUnits []models.RobboUnitCore, err error)
}

type usersByRobboUnitIdProvider interface {
	GetStudentsByRobboUnitId(robboUnitId uint) (students []models.UserCore, err error)
	GetTeachersByRobboUnitId(robboUnitId uint) (users []models.UserCore, err error)
}

type UserServiceImpl struct {
	userGateway                   gateways.UserGateway
	usersByRobboUnitIdProvider    usersByRobboUnitIdProvider
	robboUnitsByUnitAdminProvider robboUnitsByUnitAdminProvider
	parentByChildIdProvider       parentByChildIdProvider
	robboGroupRelProvider         robboGroupRelProvider
}

func (u UserServiceImpl) GetStudentsByTeacherId(teacherId uint) (students []models.UserCore, err error) {
	robboGroups, err := u.robboGroupRelProvider.GetRobboGroupsByUserId(teacherId)
	if err != nil {
		return []models.UserCore{}, err
	}
	for _, robboGroup := range robboGroups {
		studentsByGroup, _, err := u.robboGroupRelProvider.GetStudentsByRobboGroupId(0, 100, robboGroup.ID)
		if err != nil {
			return []models.UserCore{}, err
		}
		students = append(students, studentsByGroup...)
	}
	return students, nil
}

func (u UserServiceImpl) GetClientsByUnitAdminId(unitAdminId uint) (clients []models.UserCore, err error) {
	students, err := u.GetStudentsByUnitAdminId(unitAdminId)
	if err != nil {
		return []models.UserCore{}, err
	}
	userIdsMap := make(map[uint]bool)
	for _, student := range students {
		parentsByChildId, err := u.parentByChildIdProvider.GetParentsByChildId(student.ID)
		if err != nil {
			return []models.UserCore{}, err
		}

		for i := 0; i < len(parentsByChildId); i++ {
			if !userIdsMap[parentsByChildId[i].ID] {
				clients = append(clients, parentsByChildId[i])
				userIdsMap[parentsByChildId[i].ID] = true
			}
		}
	}
	return clients, nil
}

func (u UserServiceImpl) GetStudentsByUnitAdminId(unitAdminId uint) (students []models.UserCore, err error) {
	robboUnits, err := u.robboUnitsByUnitAdminProvider.GetRobboUnitsByUnitAdmin(unitAdminId)
	if err != nil {
		return []models.UserCore{}, err
	}
	for _, robboUnit := range robboUnits {
		studentsByRobboUnit, err := u.usersByRobboUnitIdProvider.GetStudentsByRobboUnitId(robboUnit.ID)
		if err != nil {
			return []models.UserCore{}, err
		}
		students = append(students, studentsByRobboUnit...)
	}
	return students, nil
}

func (u UserServiceImpl) GetTeachersByUnitAdminId(unitAdminId uint) (students []models.UserCore, err error) {
	robboUnits, err := u.robboUnitsByUnitAdminProvider.GetRobboUnitsByUnitAdmin(unitAdminId)
	if err != nil {
		return []models.UserCore{}, err
	}
	for _, robboUnit := range robboUnits {
		studentsByRobboUnit, err := u.usersByRobboUnitIdProvider.GetTeachersByRobboUnitId(robboUnit.ID)
		if err != nil {
			return []models.UserCore{}, err
		}
		students = append(students, studentsByRobboUnit...)
	}
	return students, nil
}

func (u UserServiceImpl) GetUsersByEmail(page, pageSize *int, roles []models.Role, email string, clientRole models.Role) (users []models.UserCore, countRows uint, err error) {
	switch clientRole {
	case models.RoleUnitAdmin:
		for _, role := range roles {
			if role.String() == models.RoleSuperAdmin.String() || role.String() == models.RoleUnitAdmin.String() {
				return []models.UserCore{}, 0, utils.ResponseError{
					Code:    http.StatusForbidden,
					Message: consts.ErrAccessDenied,
				}
			}
		}
	}
	offset, limit := utils.GetOffsetAndLimit(page, pageSize)
	return u.userGateway.GetUsersByEmail(offset, limit, roles, email)
}

func (u UserServiceImpl) SetIsActive(id uint, isActive bool) error {
	return u.userGateway.SetIsActive(id, isActive)
}

func (u UserServiceImpl) CreateUser(user models.UserCore, clientRole models.Role) (newUser models.UserCore, err error) {
	// TODO сразу активен? надо ли высылать код
	// checking the client role for the possibility of creating a user
	if (clientRole == models.RoleUnitAdmin && user.Role.String() == models.RoleUnitAdmin.String()) ||
		user.Role.String() == models.RoleSuperAdmin.String() {
		return models.UserCore{}, utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrAccessDenied,
		}
	}
	// TODO check valid email
	exist, err := u.userGateway.DoesExistEmail(0, user.Email)
	if err != nil {
		return models.UserCore{}, err
	}
	if exist {
		return models.UserCore{}, utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrEmailAlreadyInUse,
		}
	}
	if len(user.Password) < 6 {
		return models.UserCore{}, utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrShortPassword,
		}
	}
	passwordHash := utils.HashPassword(user.Password)
	user.Password = passwordHash
	return u.userGateway.CreateUser(user)
}

func (u UserServiceImpl) DeleteUser(id uint) error {
	// TODO maybe check userById and role
	// TODO delete projects
	// delete all rels (projectPage, project, etc.)
	return u.userGateway.DeleteUser(id)
}

func (u UserServiceImpl) UpdateUser(user models.UserCore, clientRole models.Role) (updatedUser models.UserCore, err error) {
	// TODO check client role
	// какие роли кого могут обновлять
	exist, err := u.userGateway.DoesExistEmail(user.ID, user.Email)
	if err != nil {
		return models.UserCore{}, err
	}
	// checking the client role for the possibility of updating a user
	switch clientRole {
	case models.RoleUnitAdmin:
		if user.Role.String() == models.RoleSuperAdmin.String() {
			return models.UserCore{}, utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		}
	case models.RoleTeacher:
		if user.Role.String() != models.RoleTeacher.String() {
			return models.UserCore{}, utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		}
	case models.RoleParent:
		if user.Role.String() != models.RoleParent.String() {
			return models.UserCore{}, utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		}
	case models.RoleStudent:
		if user.Role.String() != models.RoleStudent.String() {
			return models.UserCore{}, utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		}
	}
	if exist {
		return models.UserCore{}, utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrEmailAlreadyInUse,
		}
	}
	return u.userGateway.UpdateUser(user)
}

func (u UserServiceImpl) GetUserById(id uint, clientRole models.Role) (models.UserCore, error) {
	user, err := u.userGateway.GetUserById(id)
	if err != nil {
		return models.UserCore{}, err
	}
	// checking the client role for the possibility of getting a user
	switch clientRole {
	case models.RoleParent:
		if user.Role.String() != models.RoleStudent.String() && user.Role.String() != models.RoleParent.String() {
			return models.UserCore{}, utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		}
	case models.RoleTeacher:
		if user.Role.String() == models.RoleUnitAdmin.String() ||
			user.Role.String() == models.RoleSuperAdmin.String() {
			return models.UserCore{}, utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		}
	case models.RoleUnitAdmin:
		if user.Role.String() == models.RoleSuperAdmin.String() {
			return models.UserCore{}, utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrAccessDenied,
			}
		}
	}
	return user, nil
}

func (u UserServiceImpl) GetAllUsers(
	page, pageSize *int,
	isActive bool,
	roles []models.Role,
	clientId uint,
	clientRole models.Role,
) (users []models.UserCore, countRows uint, err error) {
	// checking the client role for the possibility of getting a users
	switch clientRole.String() {
	case models.RoleUnitAdmin.String():
		for _, role := range roles {
			switch role.String() {
			case models.RoleStudent.String():
				students, err := u.GetStudentsByUnitAdminId(clientId)
				if err != nil {
					return []models.UserCore{}, 0, err
				}
				users = append(users, students...)

			case models.RoleTeacher.String():
				teachers, err := u.GetTeachersByUnitAdminId(clientId)
				if err != nil {
					return []models.UserCore{}, 0, err
				}
				users = append(users, teachers...)

			case models.RoleParent.String():
				clients, err := u.GetClientsByUnitAdminId(clientId)
				if err != nil {
					return []models.UserCore{}, 0, err
				}
				users = append(users, clients...)
			}

			if role.String() == models.RoleSuperAdmin.String() || role.String() == models.RoleUnitAdmin.String() {
				return []models.UserCore{}, 0, utils.ResponseError{
					Code:    http.StatusForbidden,
					Message: consts.ErrAccessDenied,
				}
			}
			return users, uint(len(users)), nil
		}
	case models.RoleTeacher.String():
		for _, role := range roles {
			if role.String() != models.RoleStudent.String() {
				return []models.UserCore{}, 0, utils.ResponseError{
					Code:    http.StatusForbidden,
					Message: consts.ErrAccessDenied,
				}
			} else {
				students, err := u.GetStudentsByTeacherId(clientId)
				if err != nil {
					return []models.UserCore{}, 0, err
				}
				return students, uint(len(students)), nil
			}
		}
	}
	offset, limit := utils.GetOffsetAndLimit(page, pageSize)
	return u.userGateway.GetAllUsers(offset, limit, isActive, roles)
}

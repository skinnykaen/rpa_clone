package services

import (
	"errors"
	"github.com/skinnykaen/rpa_clone/internal/consts"
	"github.com/skinnykaen/rpa_clone/internal/gateways"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
)

type ProjectPageService interface {
	CreateProjectPage(authorId uint) (newProjectPage models.ProjectPageCore, err error)
	DeleteProjectPage(id, clientId uint) error
	GetAllProjectPages(page, pageSize *int, role models.Role, userId uint) (projectPages []models.ProjectPageCore, countRows uint, err error)
	UpdateProjectPage(projectPage models.ProjectPageCore) (models.ProjectPageCore, error)
	GetProjectPageById(id, clientId uint) (projectPage models.ProjectPageCore, err error)
	GetProjectsPageByAuthorId(id uint, page, pageSize *int) (projectPages []models.ProjectPageCore, countRows uint, err error)
}

type ProjectPageServiceImpl struct {
	projectGateway     gateways.ProjectGateway
	projectPageGateway gateways.ProjectPageGateway
}

func (p ProjectPageServiceImpl) GetAllProjectPages(page, pageSize *int, role models.Role, userId uint) (projectPages []models.ProjectPageCore, countRows uint, err error) {
	offset, limit := utils.GetOffsetAndLimit(page, pageSize)
	if role.String() != models.RoleSuperAdmin.String() {
		return p.projectPageGateway.GetProjectPagesByAuthorId(userId, offset, limit)
	}
	return p.projectPageGateway.GetAllProjectPages(offset, limit)
}

func (p ProjectPageServiceImpl) GetProjectsPageByAuthorId(id uint, page, pageSize *int) (projectPages []models.ProjectPageCore, countRows uint, err error) {
	offset, limit := utils.GetOffsetAndLimit(page, pageSize)
	return p.projectPageGateway.GetProjectPagesByAuthorId(id, offset, limit)
}

func (p ProjectPageServiceImpl) CreateProjectPage(authorId uint) (newProjectPage models.ProjectPageCore, err error) {
	return p.projectPageGateway.CreateProjectPage(
		models.ProjectPageCore{
			AuthorID:    authorId,
			Title:       "Untitled",
			Instruction: "",
			Notes:       "",
			IsShared:    false,
		},
		models.ProjectCore{
			AuthorID: authorId,
			Json:     consts.EmptyProjectJson,
		})
}

func (p ProjectPageServiceImpl) DeleteProjectPage(id, clientId uint) error {
	return p.projectPageGateway.DeleteProjectPage(id, clientId)
}

func (p ProjectPageServiceImpl) UpdateProjectPage(projectPage models.ProjectPageCore) (models.ProjectPageCore, error) {
	//TODO check is author?
	return p.projectPageGateway.UpdateProjectPage(projectPage)
}

func (p ProjectPageServiceImpl) GetProjectPageById(id, clientId uint) (projectPage models.ProjectPageCore, err error) {
	projectPage, err = p.projectPageGateway.GetProjectPageById(id)
	if err != nil {
		return
	}
	if projectPage.IsShared {
		return projectPage, nil
	} else {
		if projectPage.AuthorID != clientId {
			return models.ProjectPageCore{}, errors.New("access denied")
		}
	}
	return projectPage, nil
}

package services

import (
	"errors"
	"github.com/skinnykaen/rpa_clone/internal/consts"
	"github.com/skinnykaen/rpa_clone/internal/gateways"
	"github.com/skinnykaen/rpa_clone/internal/models"
)

type ProjectPageService interface {
	CreateProjectPage(authorId uint) (newProjectPage models.ProjectPageCore, err error)
	DeleteProjectPage(id, clientId uint) error
	UpdateProjectPage(projectPage models.ProjectPageCore) (models.ProjectPageCore, error)
	GetProjectPageById(id, clientId uint) (projectPage models.ProjectPageCore, err error)
	GetProjectsPageByAuthorId(id uint, page, pageSize *int) (projectPages []models.ProjectPageCore, err error)
}

type ProjectPageServiceImpl struct {
	projectGateway     gateways.ProjectGateway
	projectPageGateway gateways.ProjectPageGateway
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
	//err := p.projectGateway.DeleteProject(id, clientId)
	//if err != nil {
	//	return err
	//}
	//return p.projectGateway.DeleteProject(id, clientId)
	return p.projectPageGateway.DeleteProjectPage(id, clientId)
}

func (p ProjectPageServiceImpl) UpdateProjectPage(projectPage models.ProjectPageCore) (models.ProjectPageCore, error) {
	//TODO check is author?
	return p.projectPageGateway.UpdateProjectPage(projectPage)
}

func (p ProjectPageServiceImpl) GetProjectPageById(id, clientId uint) (projectPage models.ProjectPageCore, err error) {
	//project, err := p.projectGateway.GetProjectById(id)
	//if err != nil {
	//	return
	//}
	projectPage, err = p.projectPageGateway.GetProjectPageById(id)
	if err != nil {
		return
	}
	if projectPage.IsShared {
		//projectPage.UpdatedAt = project.UpdatedAt
		return projectPage, nil
	} else {
		if projectPage.AuthorID != clientId {
			return models.ProjectPageCore{}, errors.New("access denied")
		}
	}
	//projectPage.UpdatedAt = project.UpdatedAt
	return projectPage, nil
}

func (p ProjectPageServiceImpl) GetProjectsPageByAuthorId(id uint, page, pageSize *int) (projectPages []models.ProjectPageCore, err error) {
	//project, err := p.projectGateway.GetProjectById(id)
	//if err != nil {
	//	return models.ProjectPageCore{}, err
	//}

	//TODO implement
	panic("implement me")
	//offset, limit := utils.GetOffsetAndLimit(page, pageSize)
	//p.projectPageGateway.
}

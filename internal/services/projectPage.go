package services

import (
	"github.com/spf13/viper"
	"rpa_clone/internal/consts"
	"rpa_clone/internal/gateways"
	"rpa_clone/internal/models"
	"strconv"
)

type ProjectPageService interface {
	CreateProjectPage(authorId uint) (newProjectPage models.ProjectPageCore, err error)
	DeleteProjectPage(id uint) error
	UpdateProjectPage(projectPage models.ProjectPageCore) (models.ProjectPageCore, error)
	GetProjectPageById(id uint) (projectPage models.ProjectPageCore, err error)
	GetProjectsPageByAuthorId(id uint, page, pageSize *int) (projectPages []models.ProjectPageCore, err error)
}

type ProjectPageServiceImpl struct {
	projectGateway     gateways.ProjectGateway
	projectPageGateway gateways.ProjectPageGateway
}

func (p ProjectPageServiceImpl) CreateProjectPage(authorId uint) (newProjectPage models.ProjectPageCore, err error) {
	newProject, err := p.projectGateway.CreateProject(
		models.ProjectCore{
			AuthorID: authorId,
			Json:     consts.EmptyProjectJson,
		})
	if err != nil {
		return models.ProjectPageCore{}, err
	}
	return p.projectPageGateway.CreateProjectPage(
		models.ProjectPageCore{
			Title:       "Untitled",
			ProjectID:   newProject.ID,
			Instruction: "",
			Notes:       "",
			LinkToScratch: viper.GetString("projectPage.scratchLink") +
				"?#" + strconv.FormatUint(uint64(newProject.ID), 10),
			IsShared: false,
		})
}

func (p ProjectPageServiceImpl) DeleteProjectPage(id uint) error {
	err := p.projectGateway.DeleteProject(id)
	if err != nil {
		return err
	}
	return p.projectPageGateway.DeleteProjectPage(id)
}

func (p ProjectPageServiceImpl) UpdateProjectPage(projectPage models.ProjectPageCore) (models.ProjectPageCore, error) {
	return p.projectPageGateway.UpdateProjectPage(projectPage)
}

func (p ProjectPageServiceImpl) GetProjectPageById(id uint) (projectPage models.ProjectPageCore, err error) {
	project, err := p.projectGateway.GetProjectById(id)
	if err != nil {
		return
	}
	projectPage, err = p.projectPageGateway.GetProjectPageById(id)
	if err != nil {
		return
	}
	projectPage.UpdatedAt = project.UpdatedAt
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

package services

import (
	"github.com/skinnykaen/rpa_clone/internal/gateways"
	"github.com/skinnykaen/rpa_clone/internal/models"
)

type ProjectService interface {
	CreateProject(project models.ProjectCore) (newProject models.ProjectCore, err error)
	DeleteProject(id uint) error
	UpdateProject(project models.ProjectCore) (updatedProject models.ProjectCore, err error)
	GetProjectById(id uint) (project models.ProjectCore, err error)
	GetProjectsByAuthorId(id uint) (projects []models.ProjectCore, err error)
}

type ProjectServiceImpl struct {
	projectGateway gateways.ProjectGateway
}

func (p ProjectServiceImpl) CreateProject(project models.ProjectCore) (newProject models.ProjectCore, err error) {
	//TODO implement me
	panic("implement me")
}

func (p ProjectServiceImpl) DeleteProject(id uint) error {
	//TODO implement me
	panic("implement me")
}

func (p ProjectServiceImpl) UpdateProject(project models.ProjectCore) (updatedProject models.ProjectCore, err error) {
	//TODO implement me
	panic("implement me")
}

func (p ProjectServiceImpl) GetProjectById(id uint) (project models.ProjectCore, err error) {
	//TODO implement me
	panic("implement me")
}

func (p ProjectServiceImpl) GetProjectsByAuthorId(id uint) (projects []models.ProjectCore, err error) {
	//TODO implement me
	panic("implement me")
}

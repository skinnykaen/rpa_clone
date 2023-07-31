package services

import (
	"errors"
	"github.com/skinnykaen/rpa_clone/internal/gateways"
	"github.com/skinnykaen/rpa_clone/internal/models"
)

type ProjectService interface {
	UpdateProject(project models.ProjectCore) (updatedProject models.ProjectCore, err error)
	GetProjectById(id, clientId uint, clientRole models.Role) (project models.ProjectCore, err error)
}

type ProjectServiceImpl struct {
	projectGateway gateways.ProjectGateway
}

func (p ProjectServiceImpl) UpdateProject(project models.ProjectCore) (updatedProject models.ProjectCore, err error) {
	//TODO check is author?
	return p.projectGateway.UpdateProject(project)
}

func (p ProjectServiceImpl) GetProjectById(id, clientId uint, clientRole models.Role) (project models.ProjectCore, err error) {
	project, err = p.projectGateway.GetProjectById(id)
	if err != nil {
		return models.ProjectCore{}, err
	}
	if project.IsShared || clientRole.String() == models.RoleSuperAdmin.String() {
		return project, nil
	} else {
		if project.AuthorID != clientId {
			return models.ProjectCore{}, errors.New("access denied")
		}
	}
	return
}

package gateways

import (
	"gorm.io/gorm/clause"
	"rpa_clone/internal/db"
	"rpa_clone/internal/models"
)

type ProjectGateway interface {
	CreateProject(project models.ProjectCore) (models.ProjectCore, error)
	DeleteProject(id uint) error
	UpdateProject(project models.ProjectCore) (updatedProject models.ProjectCore, err error)
	GetProjectById(id uint) (project models.ProjectCore, err error)
	GetProjectsByAuthorId(id uint) (projects []models.ProjectCore, err error)
}

type ProjectGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (p ProjectGatewayImpl) GetProjectsByAuthorId(id uint) (projects []models.ProjectCore, err error) {
	err = p.postgresClient.Db.Where("author_id = ?", id).Find(&projects).Error
	return
}

func (p ProjectGatewayImpl) CreateProject(project models.ProjectCore) (models.ProjectCore, error) {
	result := p.postgresClient.Db.Create(project).Clauses(clause.Returning{})
	return project, result.Error
}

func (p ProjectGatewayImpl) DeleteProject(id uint) error {
	return p.postgresClient.Db.Delete(&models.ProjectCore{ID: id}).Error
}

func (p ProjectGatewayImpl) UpdateProject(project models.ProjectCore) (updatedProject models.ProjectCore, err error) {
	result := p.postgresClient.Db.Where(&models.ProjectCore{ID: project.ID}).Updates(project)
	return project, result.Error
}

func (p ProjectGatewayImpl) GetProjectById(id uint) (project models.ProjectCore, err error) {
	err = p.postgresClient.Db.First(&project, id).Error
	return
}

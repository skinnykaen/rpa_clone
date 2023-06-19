package gateways

import (
	"errors"
	"fmt"
	"github.com/skinnykaen/rpa_clone/internal/db"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"gorm.io/gorm/clause"
)

type ProjectGateway interface {
	CreateProject(project models.ProjectCore) (models.ProjectCore, error)
	DeleteProject(id, clientId uint) error
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
	result := p.postgresClient.Db.Create(&project).Clauses(clause.Returning{})
	return project, result.Error
}

func (p ProjectGatewayImpl) DeleteProject(id, clientId uint) error {
	result := p.postgresClient.Db.Where("author_id = ?", clientId).Delete(&models.ProjectCore{ID: id})
	var countRows int
	result.Row().Scan(&countRows)
	fmt.Println(countRows)
	if countRows == 0 && result.Error == nil {
		return errors.New("access denied")
	}
	return result.Error
}

func (p ProjectGatewayImpl) UpdateProject(project models.ProjectCore) (updatedProject models.ProjectCore, err error) {
	result := p.postgresClient.Db.Where(&models.ProjectCore{ID: project.ID}).Updates(project)
	return project, result.Error
}

func (p ProjectGatewayImpl) GetProjectById(id uint) (project models.ProjectCore, err error) {
	err = p.postgresClient.Db.First(&project, id).Error
	return
}

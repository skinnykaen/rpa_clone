package gateways

import (
	"github.com/skinnykaen/rpa_clone/internal/db"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"gorm.io/gorm/clause"
)

type ProjectPageGateway interface {
	CreateProjectPage(projectPage models.ProjectPageCore) (newProjectPage models.ProjectPageCore, err error)
	DeleteProjectPage(id, clientId uint) error
	UpdateProjectPage(projectPage models.ProjectPageCore) (updatedProjectPage models.ProjectPageCore, err error)
	GetProjectPageById(id uint) (projectPage models.ProjectPageCore, err error)
	SetIsShared(id uint, isShared bool) error
}

type ProjectPageGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (p ProjectPageGatewayImpl) SetIsShared(id uint, isShared bool) error {
	return p.postgresClient.Db.Where(&models.ProjectCore{}, id).Update("is_shared", isShared).Error
}

func (p ProjectPageGatewayImpl) CreateProjectPage(projectPage models.ProjectPageCore) (models.ProjectPageCore, error) {
	result := p.postgresClient.Db.Create(&projectPage).Clauses(clause.Returning{})
	return projectPage, result.Error
}

func (p ProjectPageGatewayImpl) DeleteProjectPage(id, clientId uint) error {
	return p.postgresClient.Db.Where("author_id = ?", clientId).Delete(&models.ProjectPageCore{}, id).Error
	//return p.postgresClient.Db.Model(&models.ProjectPageCore{}).Where("id = ? AND author_id = ?", id, clientId).
	//Association("Project").Unscoped().Clear()
}

func (p ProjectPageGatewayImpl) UpdateProjectPage(projectPage models.ProjectPageCore) (updatedProjectPage models.ProjectPageCore, err error) {
	err = p.postgresClient.Db.Where("id = ?", projectPage.ID).
		Updates(
			map[string]interface{}{
				"link_scratch": projectPage.LinkToScratch,
				"title":        projectPage.Title,
				"instruction":  projectPage.Instruction,
				"notes":        projectPage.Notes,
				"is_shared":    projectPage.IsShared,
			},
		).Error
	return
}

func (p ProjectPageGatewayImpl) GetProjectPageById(id uint) (projectPage models.ProjectPageCore, err error) {
	err = p.postgresClient.Db.First(&projectPage, id).Error
	return
}

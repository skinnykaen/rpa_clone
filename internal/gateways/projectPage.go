package gateways

import (
	"gorm.io/gorm/clause"
	"rpa_clone/internal/db"
	"rpa_clone/internal/models"
)

type ProjectPageGateway interface {
	CreateProjectPage(projectPage models.ProjectPageCore) (newProjectPage models.ProjectPageCore, err error)
	DeleteProjectPage(id uint) error
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

func (p ProjectPageGatewayImpl) DeleteProjectPage(id uint) error {
	return p.postgresClient.Db.Delete(&models.ProjectPageCore{}, id).Error
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

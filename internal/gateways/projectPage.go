package gateways

import (
	"github.com/skinnykaen/rpa_clone/internal/db"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
)

type ProjectPageGateway interface {
	CreateProjectPage(projectPage models.ProjectPageCore, project models.ProjectCore) (newProjectPage models.ProjectPageCore, err error)
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

func (p ProjectPageGatewayImpl) CreateProjectPage(projectPage models.ProjectPageCore, project models.ProjectCore) (models.ProjectPageCore, error) {
	if err := p.postgresClient.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&project).Clauses(clause.Returning{}).Error; err != nil {
			return err
		}
		projectPage.Project = project
		projectPage.LinkToScratch = viper.GetString("projectPage.scratchLink") +
			"?#" + strconv.FormatUint(uint64(project.ID), 10)
		if err := tx.Create(&projectPage).Clauses(clause.Returning{}).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return models.ProjectPageCore{}, err
	}
	return projectPage, nil
}

func (p ProjectPageGatewayImpl) DeleteProjectPage(id, clientId uint) error {
	//return p.postgresClient.Db.Where("author_id = ?", clientId).Delete(&models.ProjectPageCore{}, id).Error
	if err := p.postgresClient.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("author_id = ?", clientId).Delete(&models.ProjectCore{}, id).Error; err != nil {
			return err
		}
		if err := tx.Where("author_id = ?", clientId).Delete(&models.ProjectPageCore{}, id).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
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
	err = p.postgresClient.Db.Model(&models.ProjectPageCore{}).Preload("Project").First(&projectPage, id).Error
	return
}

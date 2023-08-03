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
	GetAllProjectPages(offset, limit int) (projectPages []models.ProjectPageCore, countRows uint, err error)
	GetProjectPagesByAuthorId(id uint, offset, limit int) (projectPages []models.ProjectPageCore, countRows uint, err error)
	UpdateProjectPage(projectPage models.ProjectPageCore) (updatedProjectPage models.ProjectPageCore, err error)
	GetProjectPageById(id uint) (projectPage models.ProjectPageCore, err error)
	SetIsShared(id uint, isShared bool) error
	SetIsBanned(id uint, isBanned bool) error
}

type ProjectPageGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (p ProjectPageGatewayImpl) SetIsBanned(id uint, isBanned bool) error {
	if err := p.postgresClient.Db.Transaction(func(tx *gorm.DB) error {
		if err := p.postgresClient.Db.First(&models.ProjectPageCore{ID: id}).Updates(map[string]interface{}{
			"is_banned": isBanned},
		).Error; err != nil {
			return err
		}
		if err := p.postgresClient.Db.First(&models.ProjectCore{ID: id}).Updates(map[string]interface{}{
			"is_banned": isBanned},
		).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (p ProjectPageGatewayImpl) GetProjectPagesByAuthorId(id uint, offset, limit int) (projectPages []models.ProjectPageCore, countRows uint, err error) {
	var count int64
	result := p.postgresClient.Db.Limit(limit).Offset(offset).Where("author_id = ? AND is_banned = ?", id, false).
		Find(&projectPages).Preload("Project")
	if result.Error != nil {
		return []models.ProjectPageCore{}, 0, result.Error
	}
	result.Count(&count)
	return projectPages, uint(count), result.Error
}

func (p ProjectPageGatewayImpl) GetAllProjectPages(offset, limit int) (projectPages []models.ProjectPageCore, countRows uint, err error) {
	var count int64
	result := p.postgresClient.Db.Limit(limit).Offset(offset).Find(&projectPages).Preload("Project")
	if result.Error != nil {
		return []models.ProjectPageCore{}, 0, result.Error
	}
	result.Count(&count)
	return projectPages, uint(count), result.Error
}

func (p ProjectPageGatewayImpl) SetIsShared(id uint, isShared bool) error {
	return p.postgresClient.Db.Where(&models.ProjectPageCore{}, id).Update("is_shared", isShared).Error
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
	err = p.postgresClient.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&projectPage).Clauses(clause.Returning{}).Take(&models.ProjectPageCore{}, projectPage.ID).
			Updates(
				map[string]interface{}{
					"link_to_scratch": projectPage.LinkToScratch,
					"title":           projectPage.Title,
					"instruction":     projectPage.Instruction,
					"notes":           projectPage.Notes,
					"is_shared":       projectPage.IsShared,
				},
			).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.ProjectCore{}).Clauses(clause.Returning{}).Take(&models.ProjectCore{}, projectPage.ID).
			Updates(
				map[string]interface{}{
					"is_shared": projectPage.IsShared,
				},
			).Error; err != nil {
			return err
		}
		return nil
	})
	return updatedProjectPage, nil
}

func (p ProjectPageGatewayImpl) GetProjectPageById(id uint) (projectPage models.ProjectPageCore, err error) {
	err = p.postgresClient.Db.Model(&models.ProjectPageCore{}).Preload("Project").First(&projectPage, id).Error
	return
}

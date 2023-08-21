package gateways

import (
	"github.com/skinnykaen/rpa_clone/internal/consts"
	"github.com/skinnykaen/rpa_clone/internal/db"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
)

//TODO maybe return []student or []parent а не rels

type ParentRelGateway interface {
	CreateRel(rel models.ParentRelCore) (models.ParentRelCore, error)
	DeleteRel(rel models.ParentRelCore) (err error)
	GetRelById(id uint) (models.ParentRelCore, error)
	GetChildrenByParentId(parentId uint) (students []models.UserCore, err error)
	GetParentsByChildId(childId uint) (parents []models.UserCore, err error)
}

type ParentRelGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (p ParentRelGatewayImpl) DeleteRel(rel models.ParentRelCore) error {
	if err := p.postgresClient.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("parent_id = ? AND child_id = ?", rel.ParentID, rel.ChildID).First(&rel).Error; err != nil {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: consts.ErrNotFoundInDB,
			}
		}
		if err := tx.Delete(&rel).Error; err != nil {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (p ParentRelGatewayImpl) GetRelById(id uint) (rel models.ParentRelCore, err error) {
	if err := p.postgresClient.Db.First(&rel, id).Error; err != nil {
		return models.ParentRelCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return rel, nil
}

func (p ParentRelGatewayImpl) CreateRel(rel models.ParentRelCore) (models.ParentRelCore, error) {
	if err := p.postgresClient.Db.Transaction(func(tx *gorm.DB) error {
		child := models.UserCore{ID: rel.ChildID}
		parent := models.UserCore{ID: rel.ParentID}
		if err := tx.First(&child).Error; err != nil {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		if err := tx.First(&parent).Error; err != nil {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		if child.Role != models.RoleStudent || parent.Role != models.RoleParent {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: "incorrect parent or child id",
			}
		}
		if err := tx.Create(&rel).Clauses(clause.Returning{}).Error; err != nil {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		return nil
	}); err != nil {
		return models.ParentRelCore{}, err
	}
	return rel, nil
}

func (p ParentRelGatewayImpl) GetChildrenByParentId(parentId uint) (students []models.UserCore, err error) {
	var rels []models.ParentRelCore
	if err := p.postgresClient.Db.Where(models.ParentRelCore{ParentID: parentId}).Preload("Child").Find(&rels).Error; err != nil {
		return []models.UserCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	for _, rel := range rels {
		students = append(students, rel.Child)
	}
	return students, nil
}

func (p ParentRelGatewayImpl) GetParentsByChildId(childId uint) (parents []models.UserCore, err error) {
	var rels []models.ParentRelCore
	if err := p.postgresClient.Db.Where(models.ParentRelCore{ParentID: childId}).Preload("Parent").Find(&rels).Error; err != nil {
		return []models.UserCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	for _, rel := range rels {
		parents = append(parents, rel.Parent)
	}
	return parents, nil
}

package services

import (
	"github.com/skinnykaen/rpa_clone/internal/gateways"
	"github.com/skinnykaen/rpa_clone/internal/models"
)

type ParentRelService interface {
	CreateRel(parentId, childId uint) (models.ParentRelCore, error)
	DeleteRel(parentId, childId uint) (err error)
	GetRelById(id uint) (models.ParentRelCore, error)
	GetChildrenByParentId(parentId uint) (students []models.UserCore, err error)
	GetParentsByChildId(childId uint) (parents []models.UserCore, err error)
}

type ParentRelServiceImpl struct {
	parentRelGateway gateways.ParentRelGateway
}

func (p ParentRelServiceImpl) DeleteRel(parentId, childId uint) (err error) {
	return p.parentRelGateway.DeleteRel(models.ParentRelCore{ParentID: parentId, ChildID: childId})
}

func (p ParentRelServiceImpl) GetRelById(id uint) (models.ParentRelCore, error) {
	return p.parentRelGateway.GetRelById(id)
}

func (p ParentRelServiceImpl) CreateRel(parentId, childId uint) (models.ParentRelCore, error) {
	return p.parentRelGateway.CreateRel(models.ParentRelCore{ParentID: parentId, ChildID: childId})
}

func (p ParentRelServiceImpl) GetChildrenByParentId(parentId uint) (rels []models.UserCore, err error) {
	return p.parentRelGateway.GetChildrenByParentId(parentId)
}

func (p ParentRelServiceImpl) GetParentsByChildId(childId uint) (rels []models.UserCore, err error) {
	return p.parentRelGateway.GetParentsByChildId(childId)
}

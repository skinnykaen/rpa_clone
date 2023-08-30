package gateways

import (
	"github.com/skinnykaen/rpa_clone/internal/db"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
	"net/http"
)

type MessageGateway interface {
	PostMessage(message models.MessageCore) (models.MessageCore, error)
	DeleteMessage(id uint) error
	UpdateMessage(id uint, payload string) (models.MessageCore, error)

	MessagesFromUser(receiverId, senderId uint) ([]models.MessageCore, error)

	GetMessageByID(messageID uint) (models.MessageCore, error)
}

type MessageGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (m MessageGatewayImpl) PostMessage(message models.MessageCore) (models.MessageCore, error) {

	err := m.postgresClient.Db.Create(&message).Error

	if err != nil {
		return models.MessageCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error()}
	}

	m.postgresClient.Db.Preload("Sender").Preload("Receiver").First(&message)

	if err != nil {
		return models.MessageCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error()}
	}

	return message, nil

}

func (m MessageGatewayImpl) DeleteMessage(id uint) error {

	// TODO: удалять ли безвозвратно?
	err := m.postgresClient.Db.Delete(&models.MessageCore{}, id).Error

	if err != nil {
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

func (m MessageGatewayImpl) UpdateMessage(id uint, payload string) (models.MessageCore, error) {

	var message models.MessageCore

	err := m.postgresClient.Db.Model(&message).Where("id = ?", id).
		Update("payload", payload).Error

	if err != nil {
		return models.MessageCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	m.postgresClient.Db.Preload("Sender").Preload("Receiver").First(&message)

	return message, nil
}

func (m MessageGatewayImpl) MessagesFromUser(receiverId, senderId uint) ([]models.MessageCore, error) {
	var messagesFromUser []models.MessageCore

	err := m.postgresClient.Db.Preload("Receiver").Preload("Sender").
		Where("sender_id = ? AND receiver_id = ?", senderId, receiverId).Order("id desc").Find(&messagesFromUser).Error

	if err != nil {
		return nil, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return messagesFromUser, nil
}

func (m MessageGatewayImpl) GetMessageByID(messageID uint) (models.MessageCore, error) {
	var message models.MessageCore

	err := m.postgresClient.Db.First(&message, messageID).Error

	return message, err
}

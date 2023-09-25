package gateways

import (
	"github.com/skinnykaen/rpa_clone/internal/db"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
	"gorm.io/gorm/clause"
	"net/http"
)

type MessageGateway interface {
	PostMessage(message models.MessageCore) (models.MessageCore, error)
	DeleteMessage(id uint) (receiverID uint, err error)
	UpdateMessage(id uint, payload string) (models.MessageCore, error)

	MessagesFromUser(receiverId, senderId uint) ([]models.MessageCore, error)

	GetMessageById(messageId uint) (models.MessageCore, error)
	GetMessagesByChatId(chatId uint) ([]models.MessageCore, error)
}

type MessageGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (m MessageGatewayImpl) GetMessagesByChatId(chatId uint) ([]models.MessageCore, error) {
	var messages []models.MessageCore

	if err := m.postgresClient.Db.Preload("Receiver").Preload("Sender").
		Where("chat_id = ?", chatId).Find(&messages).Error; err != nil {
		return nil, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return messages, nil
}

func (m MessageGatewayImpl) PostMessage(message models.MessageCore) (models.MessageCore, error) {

	if err := m.postgresClient.Db.Create(&message).Error; err != nil {
		return models.MessageCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error()}
	}

	if err := m.postgresClient.Db.Preload("Sender").Preload("Receiver").First(&message).Error; err != nil {
		return models.MessageCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error()}
	}

	return message, nil

}

func (m MessageGatewayImpl) DeleteMessage(id uint) (receiverID uint, err error) {

	var mes models.MessageCore

	if err := m.postgresClient.Db.Clauses(clause.Returning{Columns: []clause.Column{{Name: "receiver_id"}}}).
		Delete(&mes, id).Error; err != nil {
		return 0, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return mes.ReceiverID, nil
}

func (m MessageGatewayImpl) UpdateMessage(id uint, payload string) (models.MessageCore, error) {

	var message models.MessageCore

	if err := m.postgresClient.Db.Model(&message).Where("id = ?", id).
		Update("payload", payload).Error; err != nil {
		return models.MessageCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	if err := m.postgresClient.Db.Preload("Sender").Preload("Receiver").
		First(&message, id).Error; err != nil {
		return models.MessageCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return message, nil
}

func (m MessageGatewayImpl) MessagesFromUser(receiverId, senderId uint) ([]models.MessageCore, error) {
	var messagesFromUser []models.MessageCore

	if err := m.postgresClient.Db.Preload("Receiver").Preload("Sender").
		Where("sender_id = ? AND receiver_id = ?", senderId, receiverId).
		//Or("sender_id = ? AND receiver_id = ?", receiverId, senderId).
		Order("id desc").Find(&messagesFromUser).Error; err != nil {
		return nil, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return messagesFromUser, nil
}

func (m MessageGatewayImpl) GetMessageById(messageId uint) (models.MessageCore, error) {
	var message models.MessageCore

	if err := m.postgresClient.Db.First(&message, messageId).Error; err != nil {
		return models.MessageCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return message, nil
}

package services

import (
	"github.com/skinnykaen/rpa_clone/internal/consts"
	"github.com/skinnykaen/rpa_clone/internal/gateways"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
	"net/http"
)

type MessageService interface {
	PostMessage(message models.MessageCore, clientRole models.Role) (models.MessageCore, error)
	DeleteMessage(id, userID uint) error
	UpdateMessage(id uint, payload string, userID uint) (models.MessageCore, error)

	MessagesFromUser(receiverId, senderId uint, count int, cursor string) ([]models.MessageCore, error)
}

type MessageServiceImpl struct {
	messageGateway gateways.MessageGateway
	getterUserByID GetterUserByID
	getterChat     GetterChat
}

type GetterUserByID interface {
	GetUserById(id uint) (user models.UserCore, err error)
}

type GetterChat interface {
	CreateChat(user1ID, user2ID uint) (models.ChatCore, error)
}

func (m MessageServiceImpl) PostMessage(message models.MessageCore, clientRole models.Role) (models.MessageCore, error) {

	// Получаем роль получателя сообщения
	receiver, err := m.getterUserByID.GetUserById(message.ReceiverID)

	if err != nil {
		return models.MessageCore{}, err
	}

	// Проверяем, имеет ли отправитель доступ отправить сообщение получателю
	if err = CheckAccessForMessaging(clientRole, receiver.Role); err != nil {
		return models.MessageCore{}, err
	}

	chat, err := m.getterChat.CreateChat(message.ReceiverID, message.SenderID)

	if err != nil {
		return models.MessageCore{}, err
	}

	message.Chat = chat
	message.ChatID = chat.ID

	message, err = m.messageGateway.PostMessage(message)

	if err != nil {
		return models.MessageCore{}, err
	}

	return message, nil

}

func (m MessageServiceImpl) DeleteMessage(id, userID uint) error {
	message, err := m.messageGateway.GetMessageByID(id)

	if err != nil {
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	if message.SenderID != userID {
		return utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrAccessDenied,
		}
	}

	return m.messageGateway.DeleteMessage(id)
}

func (m MessageServiceImpl) UpdateMessage(id uint, payload string, userID uint) (models.MessageCore, error) {
	message, err := m.messageGateway.GetMessageByID(id)

	if err != nil {
		return models.MessageCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	if message.SenderID != userID {
		return models.MessageCore{}, utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrAccessDenied,
		}
	}

	return m.messageGateway.UpdateMessage(id, payload)
}

func (m MessageServiceImpl) MessagesFromUser(receiverId, senderId uint, count int, cursor string) ([]models.MessageCore, error) {
	//TODO implement me
	panic("implement me")
}

func CheckAccessForMessaging(senderRole, receiverRole models.Role) error {
	//TODO check permissions
	return nil
	/*	return utils.ResponseError{
		Code:    http.StatusForbidden,
		Message: consts.ErrAccessDenied}*/
}

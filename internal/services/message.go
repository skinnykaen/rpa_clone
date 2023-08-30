package services

import (
	"encoding/base64"
	"github.com/skinnykaen/rpa_clone/internal/consts"
	"github.com/skinnykaen/rpa_clone/internal/gateways"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
	"net/http"
	"strconv"
	"strings"
)

type MessageService interface {
	PostMessage(message models.MessageCore, clientRole models.Role) (models.MessageCore, error)
	DeleteMessage(id, userID uint) error
	UpdateMessage(id uint, payload string, userID uint) (models.MessageCore, error)

	MessagesFromUser(receiverId, senderId uint, count *int, cursor *string, userID uint) ([]models.MessageCore, int, int, error)
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

	if message.SenderID == message.ReceiverID {
		return models.MessageCore{}, utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: "sending a message to yourself"}
	}

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

func (m MessageServiceImpl) MessagesFromUser(receiverId, senderId uint, count *int, cursor *string, userID uint) (messages []models.MessageCore, from int, to int, err error) {

	if userID != receiverId && userID != senderId {
		return nil, 0, 0, utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrAccessDenied,
		}
	}

	from = 0

	if cursor != nil {
		b, err := base64.StdEncoding.DecodeString(*cursor)

		if err != nil {
			return nil, 0, 0, err
		}

		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))

		if err != nil {
			return nil, 0, 0, err
		}

		from = i
	}

	messagesFromUser, err := m.messageGateway.MessagesFromUser(receiverId, senderId)

	if err != nil {
		return nil, 0, 0, err
	}

	to = len(messagesFromUser)
	if count != nil {
		to = from + *count

		if to > len(messagesFromUser) {
			to = len(messagesFromUser)
		}
	}

	return messagesFromUser, from, to, nil
}

var Permissions = map[models.Role][]models.Role{
	models.RoleAnonymous:  []models.Role{},
	models.RoleStudent:    []models.Role{models.RoleStudent, models.RoleParent, models.RoleTeacher},
	models.RoleParent:     []models.Role{models.RoleStudent, models.RoleTeacher, models.RoleUnitAdmin},
	models.RoleTeacher:    []models.Role{models.RoleStudent, models.RoleParent, models.RoleTeacher, models.RoleUnitAdmin},
	models.RoleUnitAdmin:  []models.Role{models.RoleParent, models.RoleTeacher, models.RoleUnitAdmin},
	models.RoleSuperAdmin: models.AllRole,
}

func CheckAccessForMessaging(senderRole, receiverRole models.Role) error {

	for _, role := range Permissions[senderRole] {
		if role == receiverRole {
			return nil
		}
	}

	return utils.ResponseError{
		Code:    http.StatusForbidden,
		Message: consts.ErrAccessDenied}
}

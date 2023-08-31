package services

import (
	"github.com/skinnykaen/rpa_clone/internal/consts"
	"github.com/skinnykaen/rpa_clone/internal/gateways"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
	"net/http"
)

type ChatService interface {
	CreateChat(user1ID, user2ID uint) (models.ChatCore, error)
	DeleteChat(chatID, userID uint) error
	Chats(userID uint) ([]models.ChatCore, error)
}

type ChatServiceImpl struct {
	chatGateway gateways.ChatGateway
}

func (c ChatServiceImpl) CreateChat(user1ID, user2ID uint) (models.ChatCore, error) {

	if user1ID == user2ID {
		return models.ChatCore{}, utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: "creating chat with yourself"}
	}

	return c.chatGateway.CreateChat(user1ID, user2ID)
}

func (c ChatServiceImpl) DeleteChat(chatID, userID uint) error {

	chat, err := c.chatGateway.ChatByID(chatID)

	if err != nil {
		return err
	}

	if chat.User1ID != userID && chat.User2ID != userID {
		return utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrAccessDenied,
		}
	}

	return c.chatGateway.DeleteChat(chatID)
}

func (c ChatServiceImpl) Chats(userID uint) ([]models.ChatCore, error) {
	return c.chatGateway.Chats(userID)
}

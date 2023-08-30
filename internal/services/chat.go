package services

import (
	"github.com/skinnykaen/rpa_clone/internal/gateways"
	"github.com/skinnykaen/rpa_clone/internal/models"
)

type ChatService interface {
	CreateChat(user1ID, user2ID uint) (models.ChatCore, error)
	DeleteChat(chatID uint) error
	Chats(userID uint) ([]models.ChatCore, error)
}

type ChatServiceImpl struct {
	chatGateway gateways.ChatGateway
}

func (c ChatServiceImpl) CreateChat(user1ID, user2ID uint) (models.ChatCore, error) {
	return c.chatGateway.CreateChat(user1ID, user2ID)
}

func (c ChatServiceImpl) DeleteChat(chatID uint) error {
	//TODO implement me
	panic("implement me")
}

func (c ChatServiceImpl) Chats(userID uint) ([]models.ChatCore, error) {
	//TODO implement me
	panic("implement me")
}

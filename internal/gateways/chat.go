package gateways

import (
	"github.com/skinnykaen/rpa_clone/internal/db"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
	"gorm.io/gorm/clause"
	"net/http"
)

type ChatGateway interface {
	CreateChat(user1ID, user2ID uint) (models.ChatCore, error)
	DeleteChat(chatID uint) error
	Chats(userID *uint) ([]models.ChatCore, error)

	ChatByUsers(user1ID, user2ID uint) (models.ChatCore, error)
}

type ChatGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (c ChatGatewayImpl) CreateChat(user1ID, user2ID uint) (models.ChatCore, error) {

	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}

	var newChat models.ChatCore

	// Ищу чат, если он уже был создан, то возвращаю его
	newChat, err := c.ChatByUsers(user1ID, user2ID)

	if err == nil {
		return newChat, nil

		// Если была возвращена ошибка и она говорит о том, что запись не найдена, то создаем новую запись
	} else if err != nil && err.Error() != "record not found" {
		return models.ChatCore{}, err
	}

	newChat = models.ChatCore{
		User1ID: user1ID,
		User2ID: user2ID,
		User1:   models.UserCore{ID: user1ID},
		User2:   models.UserCore{ID: user2ID},
	}

	err = c.postgresClient.Db.Clauses(clause.Returning{}).Unscoped().Save(&newChat).Error

	if err != nil {
		return models.ChatCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "creating chat error",
		}
	}

	return newChat, err
}

func (c ChatGatewayImpl) DeleteChat(chatID uint) error {
	//TODO implement me
	panic("implement me")
}

func (c ChatGatewayImpl) Chats(userID *uint) ([]models.ChatCore, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChatGatewayImpl) ChatByUsers(user1ID, user2ID uint) (models.ChatCore, error) {

	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}

	var chat models.ChatCore

	err := c.postgresClient.Db.Where("user1_id = ? AND user2_id = ?",
		user1ID, user2ID).First(&chat).Error

	if err != nil {
		return models.ChatCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}

	}

	return chat, nil
}

package gateways

import (
	"errors"
	"github.com/skinnykaen/rpa_clone/internal/db"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
	"gorm.io/gorm"
	"net/http"
)

type ChatGateway interface {
	CreateChat(user1ID, user2ID uint) (models.ChatCore, error)
	DeleteChat(chatID uint) error
	Chats(userID uint) ([]models.ChatCore, error)

	ChatByUsers(user1ID, user2ID uint) (models.ChatCore, error)
	ChatByID(chatID uint) (models.ChatCore, error)
}

type ChatGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (c ChatGatewayImpl) CreateChat(user1ID, user2ID uint) (models.ChatCore, error) {

	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}

	var chat models.ChatCore

	// Ищу чат, если он уже был создан, то возвращаю его
	chat, err := c.ChatByUsers(user1ID, user2ID)

	if err == nil {
		return chat, nil

		// Если была возвращена ошибка и она говорит о том, что запись не найдена, то создаем новую запись
	} else if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return models.ChatCore{}, err
	}

	chat = models.ChatCore{
		User1ID: user1ID,
		User2ID: user2ID,
		User1:   models.UserCore{ID: user1ID},
		User2:   models.UserCore{ID: user2ID},
	}

	if err := c.postgresClient.Db.Omit("User1", "User2").Create(&chat).Error; err != nil {
		return models.ChatCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return chat, err
}

func (c ChatGatewayImpl) DeleteChat(chatID uint) error {

	if err := c.postgresClient.Db.Delete(&models.ChatCore{}, chatID).Error; err != nil {
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

func (c ChatGatewayImpl) Chats(userID uint) ([]models.ChatCore, error) {
	var chats []models.ChatCore

	if err := c.postgresClient.Db.Preload("User1").Preload("User2").
		Where("user1_id = ? OR user2_id = ?", userID, userID).Find(&chats).Error; err != nil {
		return nil, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return chats, nil
}

func (c ChatGatewayImpl) ChatByUsers(user1ID, user2ID uint) (models.ChatCore, error) {

	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}

	var chat models.ChatCore

	if err := c.postgresClient.Db.Where("user1_id = ? AND user2_id = ?",
		user1ID, user2ID).Preload("User1").Preload("User2").First(&chat).Error; err != nil {
		return models.ChatCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return chat, nil
}

func (c ChatGatewayImpl) ChatByID(chatID uint) (models.ChatCore, error) {
	var chat models.ChatCore

	if err := c.postgresClient.Db.Preload("User1").Preload("User2").
		First(&chat, chatID).Error; err != nil {
		return models.ChatCore{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return chat, nil
}

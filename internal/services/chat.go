package services

import "github.com/skinnykaen/rpa_clone/internal/gateways"

type ChatService interface {
}

type ChatServiceImpl struct {
	chatGateway gateways.ChatGateway
}

package models

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

type MessageCore struct {
	ID         uint `gorm:"primaryKey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Payload    string         `gorm:"not null;type:text"`
	Chat       ChatCore
	ChatID     uint
	SenderID   uint
	Sender     UserCore `gorm:"References:ID"`
	ReceiverID uint
	Receiver   UserCore `gorm:"References:ID"`
}

func (m *MessageHTTP) ToCore() MessageCore {
	id, _ := strconv.ParseUint(m.ID, 10, 64)
	chatId, _ := strconv.ParseUint(m.ChatID, 10, 64)
	senderId, _ := strconv.ParseUint(m.Sender.ID, 10, 64)
	receiverId, _ := strconv.ParseUint(m.Receiver.ID, 10, 64)

	return MessageCore{
		ID:         uint(id),
		Payload:    m.Payload,
		ChatID:     uint(chatId),
		Chat:       ChatCore{ID: uint(chatId)},
		SenderID:   uint(senderId),
		Sender:     m.Sender.ToCore(),
		ReceiverID: uint(receiverId),
		Receiver:   m.Receiver.ToCore(),
	}
}

func (m *MessageHTTP) FromCore(messageCore MessageCore) {
	m.ID = strconv.Itoa(int(messageCore.ID))
	m.Payload = messageCore.Payload
	m.ChatID = strconv.Itoa(int(messageCore.ChatID))

	var receiver UserHTTP
	receiver.FromCore(messageCore.Receiver)
	m.Receiver = &receiver

	var sender UserHTTP
	sender.FromCore(messageCore.Receiver)
	m.Sender = &sender

	m.Time = &messageCore.CreatedAt
}

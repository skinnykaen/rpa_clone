package models

import (
	"encoding/base64"
	"fmt"
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
	sender.FromCore(messageCore.Sender)
	m.Sender = &sender

	m.Time = &messageCore.CreatedAt
}

type MessageConnection struct {
	Messages []*MessageHTTP
	From     int
	To       int
}

func (u *MessageConnection) TotalCount() int {
	return len(u.Messages)
}

func (u *MessageConnection) PageInfo() PageInfo {
	return PageInfo{
		StartCursor: EncodeCursor(u.From),
		EndCursor:   EncodeCursor(u.To - 1),
		HasNextPage: u.To < len(u.Messages),
	}
}

func EncodeCursor(i int) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("cursor%d", i+1)))
}

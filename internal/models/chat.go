package models

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

type ChatCore struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Messages    []MessageCore  `gorm:"foreignKey:ChatID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User1       UserCore
	User2       UserCore
	User1ID     uint         `gorm:"index:,unique,composite:idx_users_in_chat"`
	User2ID     uint         `gorm:"index:,unique,composite:idx_users_in_chat"`
	LastMessage *MessageCore `gorm:"foreignKey:ChatID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (c *ChatHTTP) ToCore() ChatCore {
	id, _ := strconv.ParseUint(c.ID, 10, 64)
	user1ID, _ := strconv.ParseUint(c.User1.ID, 10, 64)
	user2ID, _ := strconv.ParseUint(c.User2.ID, 10, 64)

	return ChatCore{
		ID:      uint(id),
		User1:   UserCore{ID: uint(user1ID)},
		User1ID: uint(user1ID),
		User2:   UserCore{ID: uint(user2ID)},
		User2ID: uint(user2ID),
	}
}

func (c *ChatHTTP) FromCore(chatCore ChatCore) {
	c.ID = strconv.Itoa(int(chatCore.ID))

	var user1 UserHTTP
	user1.FromCore(chatCore.User1)
	c.User1 = &user1

	var user2 UserHTTP
	user2.FromCore(chatCore.User2)
	c.User2 = &user2

	if chatCore.LastMessage != nil {
		var lastMessage MessageHTTP
		lastMessage.FromCore(*chatCore.LastMessage)
		c.LastMessage = &lastMessage
	}
}

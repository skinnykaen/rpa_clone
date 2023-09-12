package resolvers

import (
	"github.com/skinnykaen/rpa_clone/internal/consts"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
	"net/http"
	"sync"
)

type ChatObserver struct {
	Channel chan *models.ChatForSubscription
	Modes   map[models.ChatMode]bool
}

func NewChatObserver(channel chan *models.ChatForSubscription) *ChatObserver {
	return &ChatObserver{
		Channel: channel,
		Modes: map[models.ChatMode]bool{
			models.ChatModeCreate: false,
			models.ChatModeDelete: false,
		},
	}
}
func (o *ChatObserver) SubscribeOnAllModes() {
	o.Modes[models.ChatModeCreate] = true
	o.Modes[models.ChatModeDelete] = true
}
func (o *ChatObserver) SubscribeOnMode(mode models.ChatMode) {
	o.Modes[mode] = true
}

type ChatObservers struct {
	ChatObservers map[uint]*ChatObserver
	Mutex         *sync.Mutex
}
type IChatObservers interface {
	CreateObserver(userID uint, mode *models.ChatMode) (<-chan *models.ChatForSubscription, error)
	DeleteObserver(userID uint) error
	NotifyObserver(userID uint, mode models.ChatMode, chat models.ChatHTTP) error
}

func (c ChatObservers) CreateObserver(userID uint, mode *models.ChatMode) (<-chan *models.ChatForSubscription, error) {

	c.Mutex.Lock()
	observer, ok := c.ChatObservers[userID]
	c.Mutex.Unlock()

	if !ok {
		channel := make(chan *models.ChatForSubscription)

		observer := NewChatObserver(channel)
		if mode == nil {
			observer.SubscribeOnAllModes()
		} else {
			observer.SubscribeOnMode(*mode)
		}

		c.Mutex.Lock()
		c.ChatObservers[userID] = observer
		c.Mutex.Unlock()

		return channel, nil

	} else {
		if mode != nil && !observer.Modes[*mode] {
			observer.SubscribeOnMode(*mode)

		} else if mode == nil {
			observer.SubscribeOnAllModes()
		}

		return observer.Channel, nil
	}
}
func (c ChatObservers) DeleteObserver(userID uint) error {
	c.Mutex.Lock()
	delete(c.ChatObservers, userID)
	c.Mutex.Unlock()

	return nil
}
func (c ChatObservers) NotifyObserver(userID uint, mode models.ChatMode, chat models.ChatHTTP) error {
	c.Mutex.Lock()
	observer, ok := c.ChatObservers[userID]
	c.Mutex.Unlock()

	if ok && observer.Modes[mode] {
		c.Mutex.Lock()

		observer.Channel <- &models.ChatForSubscription{
			ChatHTTP: &chat,
			ChatMode: mode,
		}

		c.Mutex.Unlock()

		return nil
	}

	return utils.ResponseError{
		Code:    http.StatusInternalServerError,
		Message: consts.ErrThereIsNoObservers,
	}
}

package resolvers

import (
	"github.com/skinnykaen/rpa_clone/internal/consts"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
	"net/http"
	"sync"
)

type MessageObserver struct {
	Channel chan *models.MessageForSubscription
	Modes   map[models.MessageMode]bool
}

func NewMessageObserver(channel chan *models.MessageForSubscription) *MessageObserver {
	return &MessageObserver{
		Channel: channel,
		Modes: map[models.MessageMode]bool{
			models.MessageModeCreate: false,
			models.MessageModeDelete: false,
			models.MessageModeUpdate: false,
		},
	}
}

func (o *MessageObserver) SubscribeOnAllModes() {
	o.Modes[models.MessageModeCreate] = true
	o.Modes[models.MessageModeDelete] = true
	o.Modes[models.MessageModeUpdate] = true
}
func (o *MessageObserver) SubscribeOnMode(mode models.MessageMode) {
	o.Modes[mode] = true
}

type MessageObservers struct {
	MessageObservers map[uint]*MessageObserver
	Mutex            *sync.Mutex
}
type IMessageObservers interface {
	CreateObserver(userID uint, mode *models.MessageMode) (<-chan *models.MessageForSubscription, error)
	DeleteObserver(userID uint) error
	NotifyObserver(userID uint, mode models.MessageMode, message models.MessageCore) error
}

func (m MessageObservers) CreateObserver(userID uint, mode *models.MessageMode) (<-chan *models.MessageForSubscription, error) {

	m.Mutex.Lock()
	observer, ok := m.MessageObservers[userID]
	m.Mutex.Unlock()

	// create new observer
	if !ok {
		channel := make(chan *models.MessageForSubscription)

		observer := NewMessageObserver(channel)

		if mode == nil {
			observer.SubscribeOnAllModes()
		} else {
			// add the selected mode for the subscription
			observer.SubscribeOnMode(*mode)
		}

		m.Mutex.Lock()
		m.MessageObservers[userID] = observer
		m.Mutex.Unlock()

		return channel, nil
	} else {
		// if the selected mode is not active
		if mode != nil && !observer.Modes[*mode] {
			observer.SubscribeOnMode(*mode)
		} else if mode == nil {
			observer.SubscribeOnAllModes()
		}

		return observer.Channel, nil
	}
}
func (m MessageObservers) DeleteObserver(userID uint) error {

	m.Mutex.Lock()
	delete(m.MessageObservers, userID)
	m.Mutex.Unlock()

	return nil
}
func (m MessageObservers) NotifyObserver(userID uint, mode models.MessageMode, message models.MessageHTTP) error {

	m.Mutex.Lock()
	observer, ok := m.MessageObservers[userID]
	m.Mutex.Unlock()

	if ok {

		if observer.Modes[mode] {
			m.Mutex.Lock()

			observer.Channel <- &models.MessageForSubscription{
				MessageHTTP: &message,
				MessageMode: mode,
			}

			m.Mutex.Unlock()

			return nil
		}
	}

	return utils.ResponseError{
		Code:    http.StatusInternalServerError,
		Message: consts.ErrThereIsNoObservers,
	}
}

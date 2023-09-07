package resolvers

import (
	"github.com/skinnykaen/rpa_clone/internal/services"
	"github.com/skinnykaen/rpa_clone/pkg/logger"
	"sync"
)

type Resolver struct {
	loggers            logger.Loggers
	userService        services.UserService
	authService        services.AuthService
	projectPageService services.ProjectPageService
	settingsService    services.SettingsService
	chatService        services.ChatService
	messageService     services.MessageService

	chatObservers    ChatObservers
	messageObservers MessageObservers
}

func SetupResolvers(
	loggers logger.Loggers,
	userService services.UserService,
	authService services.AuthService,
	projectPageService services.ProjectPageService,
	settingsService services.SettingsService,
	chatService services.ChatService,
	messageService services.MessageService,

) Resolver {
	return Resolver{
		loggers:            loggers,
		userService:        userService,
		authService:        authService,
		projectPageService: projectPageService,
		settingsService:    settingsService,
		chatService:        chatService,
		messageService:     messageService,
		chatObservers:      ChatObservers{ChatObservers: map[uint]*ChatObserver{}, Mutex: &sync.Mutex{}},
		messageObservers:   MessageObservers{MessageObservers: map[uint]*MessageObserver{}, Mutex: &sync.Mutex{}},
	}
}

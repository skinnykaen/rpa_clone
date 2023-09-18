package resolvers

import (
	"github.com/skinnykaen/rpa_clone/internal/services"
	"github.com/skinnykaen/rpa_clone/pkg/logger"
	"sync"
)

type Resolver struct {
	chatService        services.ChatService
	messageService     services.MessageService

	chatObservers    ChatObservers
	messageObservers MessageObservers
  
	loggers              logger.Loggers
	userService          services.UserService
	authService          services.AuthService
	projectPageService   services.ProjectPageService
	settingsService      services.SettingsService
	parentRelService     services.ParentRelService
	robboUnitService     services.RobboUnitService
	robboGroupService    services.RobboGroupService
	robboUnitRelService  services.RobboUnitRelService
	robboGroupRelService services.RobboGroupRelService
	courseService        services.CourseService
}

func SetupResolvers(
	loggers logger.Loggers,
	userService services.UserService,
	authService services.AuthService,
	projectPageService services.ProjectPageService,
	settingsService services.SettingsService,
	chatService services.ChatService,
	messageService services.MessageService,
	parentRelService services.ParentRelService,
	robboUnitService services.RobboUnitService,
	robboGroupService services.RobboGroupService,
	robboUnitRelService services.RobboUnitRelService,
	robboGroupRelService services.RobboGroupRelService,
	courseService services.CourseService,
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
    parentRelService:     parentRelService,
		robboUnitService:     robboUnitService,
		robboGroupService:    robboGroupService,
		robboUnitRelService:  robboUnitRelService,
		robboGroupRelService: robboGroupRelService,
		courseService:        courseService,
	}
}

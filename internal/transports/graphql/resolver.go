package resolvers

import (
	"github.com/skinnykaen/rpa_clone/internal/services"
	"github.com/skinnykaen/rpa_clone/pkg/logger"
)

type Resolver struct {
	loggers            logger.Loggers
	userService        services.UserService
	authService        services.AuthService
	projectPageService services.ProjectPageService
	settingsService    services.SettingsService
	parentRelService   services.ParentRelService
}

func SetupResolvers(
	loggers logger.Loggers,
	userService services.UserService,
	authService services.AuthService,
	projectPageService services.ProjectPageService,
	settingsService services.SettingsService,
	parentRelService services.ParentRelService,
) Resolver {
	return Resolver{
		loggers:            loggers,
		userService:        userService,
		authService:        authService,
		projectPageService: projectPageService,
		settingsService:    settingsService,
		parentRelService:   parentRelService,
	}
}

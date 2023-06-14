package resolvers

import (
	"rpa_clone/internal/services"
	"rpa_clone/pkg/logger"
)

type Resolver struct {
	loggers     logger.Loggers
	userService services.UserService
	authService services.AuthService
}

func SetupResolvers(
	loggers logger.Loggers,
	userService services.UserService,
	authService services.AuthService,
) Resolver {
	return Resolver{
		loggers:     loggers,
		userService: userService,
		authService: authService,
	}
}

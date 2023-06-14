package services

import (
	"go.uber.org/fx"
	"rpa_clone/internal/gateways"
	"rpa_clone/pkg/logger"
)

type Services struct {
	fx.Out
	UserService UserService
	AuthService AuthService
}

func SetupServices(
	loggers logger.Loggers,
	userGateway gateways.UserGateway,
) Services {
	return Services{
		UserService: &UserServiceImpl{
			loggers:     loggers,
			userGateway: userGateway,
		},
		AuthService: &AuthServiceImpl{
			userGateway: userGateway,
		},
	}
}

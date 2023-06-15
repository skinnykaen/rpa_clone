package services

import (
	"go.uber.org/fx"
	"rpa_clone/internal/gateways"
)

type Services struct {
	fx.Out
	UserService        UserService
	AuthService        AuthService
	ProjectService     ProjectService
	ProjectPageService ProjectPageService
}

func SetupServices(
	userGateway gateways.UserGateway,
	projectGateway gateways.ProjectGateway,
	projectPageGateway gateways.ProjectPageGateway,
) Services {
	return Services{
		UserService: &UserServiceImpl{
			userGateway: userGateway,
		},
		AuthService: &AuthServiceImpl{
			userGateway: userGateway,
		},
		ProjectService: &ProjectServiceImpl{
			projectGateway: projectGateway,
		},
		ProjectPageService: &ProjectPageServiceImpl{
			projectGateway:     projectGateway,
			projectPageGateway: projectPageGateway,
		},
	}
}

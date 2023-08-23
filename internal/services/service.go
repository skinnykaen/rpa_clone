package services

import (
	"github.com/skinnykaen/rpa_clone/internal/gateways"
	"go.uber.org/fx"
)

type Services struct {
	fx.Out
	UserService        UserService
	AuthService        AuthService
	ProjectService     ProjectService
	ProjectPageService ProjectPageService
	SettingsService    SettingsService
	ParentRelService   ParentRelService
	RobboUnitService   RobboUnitService
}

func SetupServices(
	userGateway gateways.UserGateway,
	projectGateway gateways.ProjectGateway,
	projectPageGateway gateways.ProjectPageGateway,
	settingsGateway gateways.SettingsGateway,
	parentRelGateway gateways.ParentRelGateway,
	robboUnitGateway gateways.RobboUnitGateway,
) Services {
	return Services{
		UserService: &UserServiceImpl{
			userGateway: userGateway,
		},
		AuthService: &AuthServiceImpl{
			userGateway:     userGateway,
			settingsGateway: settingsGateway,
		},
		ProjectService: &ProjectServiceImpl{
			projectGateway: projectGateway,
		},
		ProjectPageService: &ProjectPageServiceImpl{
			projectGateway:     projectGateway,
			projectPageGateway: projectPageGateway,
		},
		SettingsService: &SettingsServiceImpl{
			settingsGateway: settingsGateway,
		},
		ParentRelService: &ParentRelServiceImpl{
			parentRelGateway: parentRelGateway,
		},
		RobboUnitService: &RobboUnitServiceImpl{
			robboUnitGateway: robboUnitGateway,
		},
	}
}

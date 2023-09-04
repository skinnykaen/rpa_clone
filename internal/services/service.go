package services

import (
	"github.com/skinnykaen/rpa_clone/internal/gateways"
	"go.uber.org/fx"
)

type Services struct {
	fx.Out
	UserService          UserService
	AuthService          AuthService
	ProjectService       ProjectService
	ProjectPageService   ProjectPageService
	SettingsService      SettingsService
	ParentRelService     ParentRelService
	RobboUnitService     RobboUnitService
	RobboGroupService    RobboGroupService
	RobboUnitRelService  RobboUnitRelService
	RobboGroupRelService RobboGroupRelService
}

func SetupServices(
	userGateway gateways.UserGateway,
	projectGateway gateways.ProjectGateway,
	projectPageGateway gateways.ProjectPageGateway,
	settingsGateway gateways.SettingsGateway,
	parentRelGateway gateways.ParentRelGateway,
	robboUnitGateway gateways.RobboUnitGateway,
	robboGroupGateway gateways.RobboGroupGateway,
	robboUnitRelGateway gateways.RobboUnitRelGateway,
	robboGroupRelGateway gateways.RobboGroupRelGateway,
) Services {
	return Services{
		UserService: &UserServiceImpl{
			userGateway:                   userGateway,
			usersByRobboUnitIdProvider:    robboUnitRelGateway,
			robboUnitsByUnitAdminProvider: robboUnitRelGateway,
			parentByChildIdProvider:       parentRelGateway,
			robboGroupRelProvider:         robboGroupRelGateway,
		},
		AuthService: &AuthServiceImpl{
			userGateway:              userGateway,
			activationByLinkProvider: settingsGateway,
		},
		ProjectService: &ProjectServiceImpl{
			projectGateway: projectGateway,
		},
		ProjectPageService: &ProjectPageServiceImpl{
			projectPageGateway: projectPageGateway,
		},
		SettingsService: &SettingsServiceImpl{
			settingsGateway: settingsGateway,
		},
		ParentRelService: &ParentRelServiceImpl{
			parentRelGateway: parentRelGateway,
		},
		RobboUnitService: &RobboUnitServiceImpl{
			robboUnitGateway:              robboUnitGateway,
			robboUnitsByUnitAdminProvider: robboUnitRelGateway,
		},
		RobboGroupService: &RobboGroupServiceImpl{
			robboGroupGateway:             robboGroupGateway,
			robboUnitsByUnitAdminProvider: robboUnitRelGateway,
			robboGroupsByTeacherProvider:  robboGroupRelGateway,
		},
		RobboUnitRelService: &RobboUnitRelServiceImpl{
			robboUnitRelGateway: robboUnitRelGateway,
		},
		RobboGroupRelService: &RobboGroupRelServiceImpl{
			robboGroupRelGateway: robboGroupRelGateway,
		},
	}
}

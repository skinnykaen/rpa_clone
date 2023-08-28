package gateways

import (
	"github.com/skinnykaen/rpa_clone/internal/db"
	"go.uber.org/fx"
)

type Gateways struct {
	fx.Out
	User        UserGateway
	ParentRel   ParentRelGateway
	Project     ProjectGateway
	ProjectPage ProjectPageGateway
	Settings    SettingsGateway
	RobboUnit   RobboUnitGateway
	RobboGroup  RobboGroupGateway
}

func SetupGateways(pc db.PostgresClient) Gateways {
	return Gateways{
		User:        UserGatewayImpl{postgresClient: pc},
		ParentRel:   ParentRelGatewayImpl{pc},
		Project:     ProjectGatewayImpl{pc},
		ProjectPage: ProjectPageGatewayImpl{pc},
		Settings:    SettingsGatewayImpl{pc},
		RobboUnit:   RobboUnitGatewayImpl{pc},
		RobboGroup:  RobboGroupGatewayImpl{pc},
	}
}

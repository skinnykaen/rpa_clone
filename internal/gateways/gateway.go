package gateways

import (
	"github.com/skinnykaen/rpa_clone/internal/db"
	"go.uber.org/fx"
)

type Gateways struct {
	fx.Out
	ChatGateway    ChatGateway
	MessageGateway MessageGateway
	User          UserGateway
	ParentRel     ParentRelGateway
	Project       ProjectGateway
	ProjectPage   ProjectPageGateway
	Settings      SettingsGateway
	RobboUnit     RobboUnitGateway
	RobboGroup    RobboGroupGateway
	RobboUnitRel  RobboUnitRelGateway
	RobboGroupRel RobboGroupRelGateway
}

func SetupGateways(pc db.PostgresClient) Gateways {
	return Gateways{		
		ChatGateway:    ChatGatewayImpl{pc},
		MessageGateway: MessageGatewayImpl{pc},
		User:          UserGatewayImpl{pc},
		ParentRel:     ParentRelGatewayImpl{pc},
		Project:       ProjectGatewayImpl{pc},
		ProjectPage:   ProjectPageGatewayImpl{pc},
		Settings:      SettingsGatewayImpl{pc},
		RobboUnit:     RobboUnitGatewayImpl{pc},
		RobboGroup:    RobboGroupGatewayImpl{pc},
		RobboUnitRel:  RobboUnitRelGatewayImpl{pc},
		RobboGroupRel: RobboGroupRelGatewayImpl{pc},
	}
}

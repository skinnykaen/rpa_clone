package gateways

import (
	"go.uber.org/fx"
	"rpa_clone/internal/db"
)

type Gateways struct {
	fx.Out
	UserGateway UserGateway
	ParentRel   ParentRel
	Project     ProjectGateway
	ProjectPage ProjectPageGateway
}

func SetupGateways(pc db.PostgresClient) Gateways {
	return Gateways{
		UserGateway: UserGatewayImpl{pc},
		ParentRel:   ParentRelGatewayImpl{pc},
		Project:     ProjectGatewayImpl{pc},
		ProjectPage: ProjectPageGatewayImpl{pc},
	}
}

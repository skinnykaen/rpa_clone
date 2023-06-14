package gateways

import (
	"go.uber.org/fx"
	"rpa_clone/internal/db"
)

type Gateways struct {
	fx.Out
	UserGateway UserGateway
}

func SetupGateways(pc db.PostgresClient) Gateways {
	return Gateways{
		UserGateway: UserGatewayImpl{
			postgresClient: pc,
		},
	}
}

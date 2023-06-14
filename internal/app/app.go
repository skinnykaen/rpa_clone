package app

import (
	"go.uber.org/fx"
	"log"
	"rpa_clone/internal/configs"
	"rpa_clone/internal/consts"
	"rpa_clone/internal/db"
	"rpa_clone/internal/gateways"
	"rpa_clone/internal/server"
	"rpa_clone/internal/services"
	resolvers "rpa_clone/internal/transports/graphql"
	"rpa_clone/pkg/logger"
)

func InvokeWith(m consts.Mode, options ...fx.Option) *fx.App {
	if err := configs.Init(m); err != nil {
		log.Fatalf("%s", err.Error())
	}
	di := []fx.Option{
		fx.Provide(func() consts.Mode { return m }),
		fx.Provide(logger.InitLogger),
		fx.Provide(db.InitPostgresClient),
		fx.Provide(gateways.SetupGateways),
		fx.Provide(services.SetupServices),
		fx.Provide(resolvers.SetupResolvers),
	}
	for _, option := range options {
		di = append(di, option)
	}
	return fx.New(di...)
}

func RunApp() {
	InvokeWith(consts.Development, fx.Invoke(server.NewServer)).Run()
}

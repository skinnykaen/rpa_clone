package server

import (
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/skinnykaen/rpa_clone/graph"
	"github.com/skinnykaen/rpa_clone/internal/consts"
	"github.com/skinnykaen/rpa_clone/internal/graphql/directives"
	resolvers "github.com/skinnykaen/rpa_clone/internal/transports/graphql"
	"github.com/skinnykaen/rpa_clone/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"net/http"
)

func NewServer(
	m consts.Mode,
	lifecycle fx.Lifecycle,
	loggers logger.Loggers,
	resolver resolvers.Resolver,
) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) (err error) {
				port := viper.GetString("graphql_server_port")
				c := graph.Config{Resolvers: &resolver}
				c.Directives.HasRole = directives.HasRole(loggers.Err)
				srv := handler.NewDefaultServer(graph.NewExecutableSchema(c))
				switch m {
				case consts.Production:
					http.Handle("/query", Auth(srv, loggers.Err))
					break
				case consts.Development:
					http.Handle("/", playground.Handler("GraphQL playground", "/query"))
					http.Handle("/query", Auth(srv, loggers.Err))
					break
				}
				loggers.Info.Printf("Connect to %s:%s/ for GraphQL playground",
					viper.GetString("server_host"),
					port,
				)
				go func() {
					loggers.Err.Fatal(http.ListenAndServe(":"+port, nil))
				}()
				return
			},
			OnStop: func(context.Context) error {
				return nil
			},
		})
}

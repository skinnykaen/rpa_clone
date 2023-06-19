package server

import (
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/skinnykaen/rpa_clone/graph"
	"github.com/skinnykaen/rpa_clone/internal/consts"
	"github.com/skinnykaen/rpa_clone/internal/graphql/directives"
	resolvers "github.com/skinnykaen/rpa_clone/internal/transports/graphql"
	http2 "github.com/skinnykaen/rpa_clone/internal/transports/http"
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
	hanlers http2.Handlers,
) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) (err error) {
				port := viper.GetString("graphql_server_port")
				c := graph.Config{Resolvers: &resolver}
				c.Directives.HasRole = directives.HasRole(loggers.Err)
				mux := http.NewServeMux()
				srv := handler.NewDefaultServer(graph.NewExecutableSchema(c))
				switch m {
				case consts.Production:
					mux.Handle("/query", Auth(srv, loggers.Err))
					break
				case consts.Development:
					mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
					mux.Handle("/query", Auth(srv, loggers.Err))
					mux.Handle("/project", Auth(hanlers.ProjectHandler, loggers.Err))
					break
				}
				loggers.Info.Printf("Connect to %s:%s/ for GraphQL playground",
					viper.GetString("server_host"),
					port,
				)
				go func() {
					loggers.Err.Fatal(http.ListenAndServe(":"+port, mux))
				}()
				return
			},
			OnStop: func(context.Context) error {
				return nil
			},
		})
}

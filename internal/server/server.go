package server

import (
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"net/http"
	"rpa_clone/graph"
	"rpa_clone/internal/consts"
	resolvers "rpa_clone/internal/transports/graphql"
	"rpa_clone/pkg/logger"
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
				//c.Directives.HasRole = HasRoleDirective()
				//c.Directives.HasRole = func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []*models.Role) (res interface{}, err error) {
				//	fmt.Println(ctx.Value(consts.KeyId))
				//	fmt.Println(ctx.Value(consts.KeyRole))
				//	return next(ctx)
				//}
				srv := handler.NewDefaultServer(graph.NewExecutableSchema(c))
				switch m {
				case consts.Production:
					http.Handle("/query", Auth(srv))
					break
				case consts.Development:
					http.Handle("/", playground.Handler("GraphQL playground", "/query"))
					http.Handle("/query", Auth(srv))
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

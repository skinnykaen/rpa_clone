package directives

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"rpa_clone/internal/consts"
	"rpa_clone/internal/models"
)

func HasRole() func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []*models.Role) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []*models.Role) (interface{}, error) {
		fmt.Println(ctx.Value(consts.KeyId))
		fmt.Println(ctx.Value(consts.KeyRole))
		return next(ctx)
	}
}

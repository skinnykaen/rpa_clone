package directives

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/skinnykaen/rpa_clone/internal/consts"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
	"log"
)

func HasRole(errLogger *log.Logger) func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []*models.Role) (interface{}, error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []*models.Role) (interface{}, error) {
		clientRole := ctx.Value(consts.KeyRole)
		if !utils.DoesHaveRole(clientRole.(models.Role), roles) {
			errLogger.Printf("%s", "access denied")
			return nil, errors.New("access denied")
		}
		return next(ctx)
	}
}

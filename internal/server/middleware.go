package server

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/spf13/viper"
	"net/http"
	"rpa_clone/internal/consts"
	"rpa_clone/internal/models"
	"rpa_clone/internal/services"
	"strings"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := w.Header().Get("Authorization")
		fmt.Println(authHeader)
		if authHeader == "" {
			r = r.WithContext(context.WithValue(r.Context(), consts.KeyId, "0"))
			r = r.WithContext(context.WithValue(r.Context(), consts.KeyRole, models.RoleAnonymous))
			next.ServeHTTP(w, r)
			return
		}
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			http.Error(w, "invalid authorization header format", http.StatusBadRequest)
			return
		}
		data, err := jwt.ParseWithClaims(headerParts[1], &services.UserClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(viper.GetString("auth_access_signing_key")), nil
			})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		claims, ok := data.Claims.(services.UserClaims)
		if !ok {
			http.Error(w, "token claims are not of type *StandardClaims", http.StatusUnauthorized)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "id", claims.Id))
		r = r.WithContext(context.WithValue(r.Context(), "role", claims.Role))
		next.ServeHTTP(w, r)
	})
}

//func HasRoleDirective() func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []*models.Role) (interface{}, error) {
//	//return func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []*models.Role) (interface{}, error) {
//	//	fmt.Println(ctx.Value(consts.KeyId))
//	//	fmt.Println(ctx.Value(consts.KeyRole))
//	//	return next(ctx)
//	//}
//}

type middleware func(http.HandlerFunc) http.HandlerFunc

func chainMiddleware(mw ...middleware) middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}
			last(w, r)
		}
	}
}

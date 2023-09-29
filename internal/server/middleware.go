package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/skinnykaen/rpa_clone/internal/consts"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/internal/services"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"strings"
	"time"
)

func WebSocketInit(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {
	// Get the token from payload
	authHeader := initPayload[consts.AuthPayload]
	fmt.Println(authHeader)
	authString, ok := authHeader.(string)
	if !ok || authString == "" {
		ctx = context.WithValue(ctx, consts.KeyId, uint(0))
		ctx = context.WithValue(ctx, consts.KeyRole, models.RoleAnonymous)
		return ctx, nil
	}

	if err := validateAuthHeader(authString); err != nil {
		return nil, err
	}

	headerParts := strings.Split(authString, " ")
	// token verification and authentication
	userId, userRole, err := getUserFromAuthentication(headerParts[1])
	if err != nil {
		return nil, err
	}

	// put it in context
	ctx = context.WithValue(ctx, consts.KeyId, userId)
	ctx = context.WithValue(ctx, consts.KeyRole, userRole)

	return ctx, nil
}

func Auth(next http.Handler, errLogger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(consts.AuthHeader)

		if authHeader == "" {
			r = r.WithContext(context.WithValue(r.Context(), consts.KeyId, "0"))
			r = r.WithContext(context.WithValue(r.Context(), consts.KeyRole, models.RoleAnonymous))
			next.ServeHTTP(w, r)
			return
		}

		if err := validateAuthHeader(authHeader); err != nil {
			errLogger.Printf("%s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		userId, userRole, err := getUserFromAuthentication(headerParts[1])

		if err != nil {
			errLogger.Printf("%s", err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), consts.KeyId, userId))
		r = r.WithContext(context.WithValue(r.Context(), consts.KeyRole, userRole))

		next.ServeHTTP(w, r)
	})
}

func getUserFromAuthentication(token string) (id uint, role models.Role, err error) {
	data, err := jwt.ParseWithClaims(token, &services.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("auth_access_signing_key")), nil
		})
	if data == nil {
		return 0, models.RoleAnonymous, errors.New(consts.ErrEmptyDataWithClaims)
	}

	claims, ok := data.Claims.(*services.UserClaims)
	if !ok {
		return 0, models.RoleAnonymous, errors.New(consts.ErrNotStandardToken)
	}
	if err != nil {
		if claims.ExpiresAt.Unix() < time.Now().Unix() {
			return 0, models.RoleAnonymous, errors.New(consts.ErrTokenExpired)
		}
		return 0, models.RoleAnonymous, err
	}

	return claims.Id, claims.Role, nil
}

func validateAuthHeader(authHeader string) error {
	// headerParts should be = ["Bearer", "<accessToken>"]
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return errors.New("invalid authorization header format")
	}
	if headerParts[0] != "Bearer" {
		return errors.New("invalid authorization header format")
	}

	return nil
}

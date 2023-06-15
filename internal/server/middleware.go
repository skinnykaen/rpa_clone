package server

import (
	"context"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"rpa_clone/internal/consts"
	"rpa_clone/internal/models"
	"rpa_clone/internal/services"
	"strings"
)

func Auth(next http.Handler, errLogger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			r = r.WithContext(context.WithValue(r.Context(), consts.KeyId, "0"))
			r = r.WithContext(context.WithValue(r.Context(), consts.KeyRole, models.RoleAnonymous))
			next.ServeHTTP(w, r)
			return
		}
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			errLogger.Printf("%s", "invalid authorization header format")
			http.Error(w, "invalid authorization header format", http.StatusBadRequest)
			return
		}
		data, err := jwt.ParseWithClaims(headerParts[1], &services.UserClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(viper.GetString("auth_access_signing_key")), nil
			})
		if err != nil {
			errLogger.Printf("%s", err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		claims, ok := data.Claims.(*services.UserClaims)
		if !ok {
			errLogger.Printf("%s", "token claims are not of type *StandardClaims")
			http.Error(w, "token claims are not of type *StandardClaims", http.StatusUnauthorized)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), consts.KeyId, claims.Id))
		r = r.WithContext(context.WithValue(r.Context(), consts.KeyRole, claims.Role))
		next.ServeHTTP(w, r)
	})
}

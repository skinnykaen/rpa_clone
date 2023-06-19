package services

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/skinnykaen/rpa_clone/internal/gateways"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
	"github.com/spf13/viper"
	"time"
)

type Tokens struct {
	Access  string
	Refresh string
}

type UserClaims struct {
	jwt.StandardClaims
	Id   uint
	Role models.Role
}

type AuthService interface {
	SignUp(newUser models.UserCore) error
	SignIn(email, password string) (Tokens, error)
	Refresh(token string) (string, error)
}

type AuthServiceImpl struct {
	userGateway gateways.UserGateway
}

func (a AuthServiceImpl) Refresh(token string) (string, error) {
	claims, err := parseToken(token, []byte(viper.GetString("auth_refresh_signing_key")))
	if err != nil {
		return "", err
	}
	user := models.UserCore{
		ID:   claims.Id,
		Role: claims.Role,
	}
	newAccessToken, err := generateToken(user, viper.GetDuration("auth_access_token_ttl"), []byte(viper.GetString("auth_access_signing_key")))
	if err != nil {
		return "", err
	}
	return newAccessToken, nil
}

func (a AuthServiceImpl) SignIn(email, password string) (Tokens, error) {
	user, err := a.userGateway.GetUserByEmail(email)
	if err = utils.ComparePassword(user.Password, password); err != nil {
		return Tokens{}, err
	}
	if err != nil {
		return Tokens{}, err
	}
	if !user.IsActive {
		return Tokens{}, errors.New("user is not active. please check your email")
	}
	access, err := generateToken(user, viper.GetDuration("auth_access_token_ttl"), []byte(viper.GetString("auth_access_signing_key")))
	if err != nil {
		return Tokens{}, err
	}
	refresh, err := generateToken(user, viper.GetDuration("auth_refresh_token_ttl"), []byte(viper.GetString("auth_refresh_signing_key")))
	if err != nil {
		return Tokens{}, err
	}
	return Tokens{Access: access, Refresh: refresh}, nil
}

func (a AuthServiceImpl) SignUp(newUser models.UserCore) error {
	if !utils.IsValidEmail(newUser.Email) {
		return errors.New("not valid email")
	}
	exist, err := a.userGateway.DoesExistEmail(0, newUser.Email)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("email already in use")
	}
	if len(newUser.Password) < 6 {
		return errors.New("please input password, at least 6 symbols")
	}

	passwordHash := utils.HashPassword(newUser.Password)
	newUser.Password = passwordHash
	newUser, err = a.userGateway.CreateUser(newUser)
	if err != nil {
		return err
	}

	//TODO config path for activation + check setting activation by code
	subject := "Ваш код активации аккаунта"
	body := "<p> Введите этот код " + fmt.Sprintf("%d", newUser.ActivationCode) +
		" для активации вашего аккаунта перейдя по ссылке http://localhost:3000/activation</p>"
	err = utils.SendEmail(subject, newUser.Email, body)
	return err
}

func generateToken(user models.UserCore, duration time.Duration, signingKey []byte) (token string, err error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(duration * time.Second)),
		},
		Id:   user.ID,
		Role: user.Role,
	}
	ss := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = ss.SignedString(signingKey)
	return
}

func parseToken(token string, key []byte) (UserClaims, error) {
	data, err := jwt.ParseWithClaims(token, UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return key, nil
		})
	if err != nil {
		return UserClaims{}, err
	}
	claims, ok := data.Claims.(UserClaims)
	if !ok {
		return UserClaims{}, errors.New("token claims are not of type *StandardClaims")
	}
	return claims, nil
}

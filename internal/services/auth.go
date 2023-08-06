package services

import (
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/skinnykaen/rpa_clone/internal/consts"
	"github.com/skinnykaen/rpa_clone/internal/gateways"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
	"github.com/spf13/viper"
	"net/http"
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
	ConfirmActivation(link string) (Tokens, error)
}

type AuthServiceImpl struct {
	userGateway     gateways.UserGateway
	settingsGateway gateways.SettingsGateway
}

func (a AuthServiceImpl) ConfirmActivation(link string) (Tokens, error) {
	activationByLink, err := a.settingsGateway.GetActivationByLink()
	if err != nil {
		return Tokens{Access: "", Refresh: ""}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	if !activationByLink {
		return Tokens{Access: "", Refresh: ""}, utils.ResponseError{
			Code:    http.StatusServiceUnavailable,
			Message: consts.ErrActivationLinkUnavailable,
		}
	}
	user, err := a.userGateway.GetUserByActivationLink(link)
	if err != nil {
		return Tokens{Access: "", Refresh: ""}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	if err := a.userGateway.SetIsActive(user.ID, true); err != nil {
		return Tokens{Access: "", Refresh: ""}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	access, err := generateToken(user, viper.GetDuration("auth_access_token_ttl"), []byte(viper.GetString("auth_access_signing_key")))
	if err != nil {
		return Tokens{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	refresh, err := generateToken(user, viper.GetDuration("auth_refresh_token_ttl"), []byte(viper.GetString("auth_refresh_signing_key")))
	if err != nil {
		return Tokens{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return Tokens{Access: access, Refresh: refresh}, nil
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
		return "", utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return newAccessToken, nil
}

func (a AuthServiceImpl) SignIn(email, password string) (Tokens, error) {
	user, err := a.userGateway.GetUserByEmail(email)
	if err != nil {
		return Tokens{}, err
	}
	if err = utils.ComparePassword(user.Password, password); err != nil {
		return Tokens{}, utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrIncorrectPasswordOrEmail,
		}
	}
	if !user.IsActive {
		return Tokens{},
			utils.ResponseError{
				Code:    http.StatusForbidden,
				Message: consts.ErrUserIsNotActive,
			}
	}
	access, err := generateToken(user, viper.GetDuration("auth_access_token_ttl"), []byte(viper.GetString("auth_access_signing_key")))
	if err != nil {
		return Tokens{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	refresh, err := generateToken(user, viper.GetDuration("auth_refresh_token_ttl"), []byte(viper.GetString("auth_refresh_signing_key")))
	if err != nil {
		return Tokens{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return Tokens{Access: access, Refresh: refresh}, nil
}

func (a AuthServiceImpl) SignUp(newUser models.UserCore) error {
	if !utils.IsValidEmail(newUser.Email) {
		return utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrIncorrectPasswordOrEmail,
		}
	}
	exist, err := a.userGateway.DoesExistEmail(0, newUser.Email)
	if err != nil {
		return err
	}
	if exist {
		return utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrIncorrectPasswordOrEmail,
		}
	}
	if len(newUser.Password) < 6 {
		return utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: consts.ErrShortPassword,
		}
	}

	passwordHash := utils.HashPassword(newUser.Password)
	newUser.Password = passwordHash
	newUser, err = a.userGateway.CreateUser(newUser)
	if err != nil {
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	activationByLink, err := a.settingsGateway.GetActivationByLink()
	if err != nil {
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	var subject, body string
	//TODO config path for activation
	if activationByLink {
		subject = "Ваша ссылка активации аккаунта"
		body = "<p>Перейдите по ссылке http://92.255.79.9/activation/" + newUser.ActivationLink +
			" для активации вашего аккаунта.</p>"
	} else {
		subject = "Активация аккаунта"
		body = "<p>На данный момент активация по ссылке недоступна. Ждите активации от администратора.</p>"
	}
	if err := utils.SendEmail(subject, newUser.Email, body); err != nil {
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
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
	return token, err
}

func parseToken(token string, key []byte) (*UserClaims, error) {
	data, err := jwt.ParseWithClaims(token, &UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return key, nil
		})
	claims, ok := data.Claims.(*UserClaims)
	if err != nil {
		if claims.ExpiresAt.Unix() < time.Now().Unix() {
			return &UserClaims{}, utils.ResponseError{
				Code:    http.StatusUnauthorized,
				Message: consts.ErrTokenExpired,
			}
		}
		return &UserClaims{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	if !ok {
		return &UserClaims{}, utils.ResponseError{
			Code:    http.StatusUnauthorized,
			Message: consts.ErrNotStandardToken,
		}
	}
	return claims, nil
}

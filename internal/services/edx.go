package services

import (
	"encoding/json"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"net/url"
	"time"
)

type EdxService interface {
	GetAllCourses() ([]byte, error)
	GetCourseById(id string) ([]byte, error)
	GetWithAuth(url string) ([]byte, error)
	RefreshToken() error
}

type EdxServiceImpl struct {
}

type NewToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func (e EdxServiceImpl) RefreshToken() error {
	if viper.GetInt64("api.token_expiration_time") < time.Now().Unix() {
		urlAddr := viper.GetString("api_urls.refreshToken")

		response, err := http.PostForm(urlAddr, url.Values{
			"grant_type":    {"client_credentials"},
			"client_id":     {viper.GetString("api.client_id")},
			"client_secret": {viper.GetString("api.client_secret")},
		})
		defer response.Body.Close()

		if err != nil {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}

		//if response.StatusCode != http.StatusOK {
		//	return utils.ResponseError{
		//		Code:    http.StatusInternalServerError,
		//		Message: "response error",
		//	}
		//}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}

		var newToken NewToken
		if err := json.Unmarshal(body, &newToken); err != nil {
			return utils.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}

		expirationTime := time.Now().Unix() + int64(newToken.ExpiresIn)
		viper.Set("api.token_expiration_time", expirationTime)
		viper.Set("api.token", newToken.AccessToken)

		return nil
	} else {
		return nil
	}
}

func (e EdxServiceImpl) GetWithAuth(url string) ([]byte, error) {
	//if err := e.RefreshToken(); err != nil {
	//	return nil, err
	//}
	//
	//bearer := "Bearer " + viper.GetString("api_token")

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	//request.Header.Add("Authorization", bearer)

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return nil, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return body, nil
}

func (e EdxServiceImpl) GetCourseById(id string) ([]byte, error) {
	request, err := http.NewRequest("GET", viper.GetString("api_urls.getCourse")+id, nil)
	if err != nil {
		return nil, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return body, nil
}

func (e EdxServiceImpl) GetAllCourses() ([]byte, error) {
	response, err := http.Get(viper.GetString("api_urls.getCourses"))
	defer response.Body.Close()

	if err != nil {
		return nil, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return body, nil
}

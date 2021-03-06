package rest

import (
	"encoding/json"
	"github.com/micro-gis/oauth-api/src/domain/users"
	errors "github.com/micro-gis/utils/rest_errors"
	"github.com/yossefaz/go-http-client/gohttp"
	"time"
)

var (
	Timeout         = 100 * time.Millisecond
	BaseURL         = "http://127.0.0.1:8086"
	usersRestClient gohttp.Client
)

func init() {
	usersRestClient = gohttp.NewBuilder().SetConnectionTimeout(Timeout).Build()
}

type RestUsersRepository interface {
	LoginUser(string, string) (*users.User, errors.RestErr)
}

type userRepository struct{}

func NewRestUsersRepository() RestUsersRepository {
	return &userRepository{}
}

func (r *userRepository) LoginUser(email string, password string) (*users.User, errors.RestErr) {
	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}
	response, err := usersRestClient.Post(BaseURL+"/users/login", request)
	if response == nil || response.StatusCode < 100 {
		return nil, errors.NewInternalServerError("invalid restClient response when trying to login user", err)
	}

	if response.StatusCode > 299 {
		restErr, err :=errors.NewRestErrorFromBytes(response.Bytes())
		if err != nil {
			return nil, errors.NewInternalServerError("invalid error interface when trying to login user", err)
		}
		return nil, restErr
	}

	var user users.User
	if err := json.Unmarshal(response.Bytes(), &user); err != nil {
		return nil, errors.NewInternalServerError("error when trying to unmarshall user response", err)
	}
	return &user, nil
}

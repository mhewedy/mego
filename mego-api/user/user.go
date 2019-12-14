package user

import (
	"encoding/json"
	"errors"
	"net/http"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type tw struct {
	IDToken string `json:"id_token"`
}

func Login(w http.ResponseWriter, r *http.Request) (interface{}, error) {

	u, err := validateAndParseLoginInput(r)
	if err != nil {
		return nil, err
	}

	token, ok := login(u)
	if !ok {
		return nil, errors.New("invalid username/password")
	}

	return &tw{IDToken: string(token)}, nil
}

func Logout(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	logout(&user{}) // TODO get user from context
	return nil, nil
}

func validateAndParseLoginInput(r *http.Request) (*user, error) {
	var i user
	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		return nil, err
	}

	if len(i.Username) == 0 || len(i.Password) == 0 {
		return nil, errors.New("username and password should not be empty")
	}
	return &i, nil
}

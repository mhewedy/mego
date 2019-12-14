package user

import (
	"github.com/mhewedy/ews"
	"github.com/mhewedy/mego/commons"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

type mockEWSClient struct {
}

func (m mockEWSClient) SendAndReceive(body []byte) ([]byte, error) {
	if strings.Contains(string(body), "<m:Email>unauthenticatedUser</m:Email>") {
		return nil, &ews.HTTPError{
			Status:     "",
			StatusCode: http.StatusUnauthorized,
		}
	}

	if strings.Contains(string(body), "<m:Email>goodUser</m:Email>") {
		return nil, nil
	}

	if strings.Contains(string(body), "<m:Email>500</m:Email>") {
		return nil, &ews.SoapError{Fault: &ews.Fault{
			Faultcode:   "ERR",
			Faultstring: "Error",
		}}
	}

	return nil, nil
}

func (m mockEWSClient) GetEWSAddr() string {
	return ""
}

func (m mockEWSClient) GetUsername() string {
	return ""
}

func Test_loginInCaseOf401(t *testing.T) {

	commons.DefaultEWSClient = &mockEWSClient{}

	_, login := login(&user{
		Username: "unauthenticatedUser",
		Password: "efg",
	})

	assert.Equal(t, false, login)
}

func Test_loginInCaseOf200(t *testing.T) {

	commons.DefaultEWSClient = &mockEWSClient{}

	token, login := login(&user{
		Username: "goodUser",
		Password: "efg",
	})

	assert.Equal(t, true, login)
	assert.NotEmpty(t, token)
}

func Test_loginInCaseOf500(t *testing.T) {

	commons.DefaultEWSClient = &mockEWSClient{}

	_, login := login(&user{
		Username: "500",
		Password: "efg",
	})

	assert.Equal(t, true, login)
}

func Test_GetUser(t *testing.T) {

	commons.DefaultEWSClient = &mockEWSClient{}
	usersDB = make(map[string]string)

	token, login := login(&user{
		Username: "goodUser",
		Password: "hisPassword",
	})
	assert.Equal(t, true, login)
	assert.NotEmpty(t, token)

	user, err := getUser(token)

	assert.NoError(t, err)
	assert.Len(t, usersDB, 1)
	assert.Equal(t, "goodUser", user.Username)
	assert.Equal(t, "hisPassword", user.Password)
}

func Test_GetUserWithInvalidUsername(t *testing.T) {

	commons.DefaultEWSClient = &mockEWSClient{}
	usersDB = make(map[string]string)

	token, login := login(&user{
		Username: "goodUser",
		Password: "hisPassword",
	})
	assert.Equal(t, true, login)
	assert.NotEmpty(t, token)
	assert.Len(t, usersDB, 1)

	user, err := getUser("invalid token")
	assert.Error(t, err)
	assert.Nil(t, user)
}

func Test_Logout(t *testing.T) {

	commons.DefaultEWSClient = &mockEWSClient{}
	usersDB = make(map[string]string)

	token, login := login(&user{
		Username: "goodUser",
		Password: "hisPassword",
	})
	assert.Equal(t, true, login)
	assert.NotEmpty(t, token)

	logout(&user{
		Username: "goodUser",
		Password: "hisPassword",
	})

	assert.Len(t, usersDB, 0)
}

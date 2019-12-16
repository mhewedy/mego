package attendess

import (
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/mhewedy/mego/commons"
	"github.com/mhewedy/mego/user"
	"net/http"
	"sync"
)

var attendOnce sync.Once

func List(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	u := context.Get(r, user.KEY).(*user.User)
	attendOnce.Do(func() {
		indexAttendees(u)
	})

	return attendeesIndex, nil
}

func Search(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	u := context.Get(r, user.KEY).(*user.User)
	attendOnce.Do(func() {
		indexAttendees(u)
	})

	var exclude []string
	err := commons.JSONDecode(r.Body, &exclude)
	if err != nil {
		return nil, err
	}

	attendees := searchAttendees(r.URL.Query().Get("q"), exclude)
	return attendees, nil
}

func GetPhoto(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	u := context.Get(r, user.KEY).(*user.User)
	ewsClient := commons.NewEWSClient(u.Username, u.Password)

	email := mux.Vars(r)["email"]
	base64, _ := getAttendeePhoto(ewsClient, email)

	return &struct {
		Base64 string `json:"base64"`
	}{Base64: base64}, nil
}

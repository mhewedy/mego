package attendess

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mhewedy/ews"
	"net/http"
	"sync"
)

var EWSClient ews.Client
var attendOnce sync.Once

func ListAttendees(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	attendOnce.Do(indexAttendees)

	return attendeesIndex, nil
}

func SearchAttendees(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	attendOnce.Do(indexAttendees)

	var exclude []string
	err := json.NewDecoder(r.Body).Decode(&exclude)
	if err != nil {
		return nil, err
	}

	attendees := searchAttendees(r.URL.Query().Get("q"), exclude)
	return attendees, nil
}

func GetPhoto(w http.ResponseWriter, r *http.Request) (interface{}, error) {

	email := mux.Vars(r)["email"]
	base64, _ := getAttendeePhoto(EWSClient, email)

	return struct {
		Base64 string `json:"base64"`
	}{Base64: base64}, nil
}

package attendess

import (
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

	attendees := searchAttendees(r.URL.Query().Get("q"))
	return attendees, nil
}

func GetPhoto(w http.ResponseWriter, r *http.Request) (interface{}, error) {

	email := mux.Vars(r)["email"]
	base64, _ := getAttendeePhoto(EWSClient, email)

	return struct {
		Base64 string `json:"base64"`
	}{Base64: base64}, nil
}

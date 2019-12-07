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

func ListAttendees(w http.ResponseWriter, r *http.Request) {
	attendOnce.Do(indexAttendees)

	json.NewEncoder(w).Encode(attendeesIndex)
}

func SearchAttendees(w http.ResponseWriter, r *http.Request) {
	attendOnce.Do(indexAttendees)

	attendees := searchAttendees(r.URL.Query().Get("q"))
	json.NewEncoder(w).Encode(attendees)
}

func GetPhoto(w http.ResponseWriter, r *http.Request) {

	email := mux.Vars(r)["email"]
	base64, err := getAttendeePhoto(EWSClient, email)

	if err != nil {
		//api.HandleError(w, err, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(struct {
		Base64 string `json:"base64"`
	}{Base64: base64})
}

package attendess

import (
	"encoding/json"
	"github.com/mhewedy/ews"
	"net/http"
	"sync"
)

var EWSClient ews.Client
var attendOnce sync.Once

func ListAttendees(w http.ResponseWriter, r *http.Request) {
	attendOnce.Do(indexAttendees)
	_ = json.NewEncoder(w).Encode(attendeesIndex)
}

func SearchAttendees(w http.ResponseWriter, r *http.Request) {
	attendOnce.Do(indexAttendees)

	_, _ = w.Write([]byte("SUCCESS"))
}

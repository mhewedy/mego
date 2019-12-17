package attendess

import (
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/mhewedy/go-conf"
	"github.com/mhewedy/mego/commons"
	"github.com/mhewedy/mego/user"
	"math"
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

	sc := conf.GetInt("attendees.search_count", 20)
	max := int(math.Min(float64(sc), float64(len(attendees))))

	return attendees[:max], nil
}

func GetByEmail(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	u := context.Get(r, user.KEY).(*user.User)
	ewsClient := commons.NewEWSClient(u.Username, u.Password)

	email := mux.Vars(r)["email"]
	return getAttendeeDetails(ewsClient, email)
}

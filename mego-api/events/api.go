package events

import (
	"github.com/gorilla/context"
	"github.com/mhewedy/ews"
	"github.com/mhewedy/ews/ewsutil"
	"github.com/mhewedy/go-conf"
	"github.com/mhewedy/mego/user"
	"net/http"
	"time"
)

type searchInput struct {
	Emails []string  `json:"emails"`
	Rooms  []string  `json:"rooms"`
	From   time.Time `json:"from"`
	To     time.Time `json:"to"`
}

type createInput struct {
	To       []string  `json:"to"`
	Optional []string  `json:"optional"`
	Subject  string    `json:"subject"`
	Body     string    `json:"body"`
	Room     string    `json:"room"`
	From     time.Time `json:"from"`
	Duration int       `json:"duration"`
}

type roomEvents struct {
	Room        string             `json:"room"`
	RoomName    string             `json:"room_name"`
	BusyDetails map[string][]event `json:"busy_details"`
	Error       string             `json:"error"`
}

type event struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	BusyType string    `json:"busy_type,omitempty"`
}

var roomIndex = 0

func Search(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	u := context.Get(r, user.KEY).(*user.User)

	input, err := parseAndValidateSearchInput(r)
	if err != nil {
		return nil, err
	}

	eventUsers := buildEventUserSlices(input, u)
	events := doSearch(eventUsers, input.From, input.To, u)

	return events, nil
}

func Create(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	u := context.Get(r, user.KEY).(*user.User)

	input, err := parseAndValidateCreateInput(r)
	if err != nil {
		return nil, err
	}

	err = doCreate(input, u)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func buildEventUserSlices(i *searchInput, u *user.User) [][]ewsutil.EventUser {

	i.From = i.From.Truncate(1 * time.Minute)
	i.To = i.To.Truncate(1 * time.Minute)

	myUsername := u.Username
	if dns := conf.Get("ews.dns_name", ""); len(dns) > 0 {
		myUsername = myUsername + "@" + dns
	}

	me := ewsutil.EventUser{
		Email:        myUsername,
		AttendeeType: ews.AttendeeTypeOrganizer,
	}
	emails := make([]ewsutil.EventUser, len(i.Emails))
	for i, ee := range i.Emails {
		emails[i] = ewsutil.EventUser{
			Email:        ee,
			AttendeeType: ews.AttendeeTypeRequired,
		}
	}
	events := make([][]ewsutil.EventUser, len(i.Rooms))

	for i, rr := range i.Rooms {
		events[i] = make([]ewsutil.EventUser, 1)
		events[i][roomIndex] = ewsutil.EventUser{
			Email:        rr,
			AttendeeType: ews.AttendeeTypeResource,
		}
		events[i] = append(events[i], emails...)
		events[i] = append(events[i], me)
	}

	return events
}

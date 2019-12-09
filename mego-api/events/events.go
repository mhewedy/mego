package events

import (
	"encoding/json"
	"github.com/mhewedy/ews"
	"github.com/mhewedy/ews/ewsutil"
	"github.com/mhewedy/mego/commons"
	"net/http"
	"time"
)

type input struct {
	Emails   []string  `json:"emails"`
	Rooms    []string  `json:"rooms"`
	From     time.Time `json:"from"`
	Duration int       `json:"duration"`
}

type roomEvents struct {
	Room   string          `json:"room"`
	Events []ewsutil.Event `json:"events"`
	Error  string          `json:"error"`
}

var EWSClient ews.Client
var roomIndex = 0

func Search(w http.ResponseWriter, r *http.Request) (interface{}, error) {

	input, err := parseAndValidate(r)
	if err != nil {
		return nil, err
	}

	eventUsers := buildEventUserSlices(input)
	events := doSearch(eventUsers, input.From, time.Duration(input.Duration)*time.Minute)

	return events, nil
}

func parseAndValidate(r *http.Request) (*input, error) {
	var i input
	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		return nil, err
	}
	if len(i.Emails) == 0 {
		return nil, commons.NewClientError("empty emails")
	}
	if len(i.Rooms) == 0 {
		return nil, commons.NewClientError("empty rooms")
	}
	if i.Duration == 0 || i.Duration%30 != 0 {
		return nil, commons.NewClientError("duration should be multiple of 30")
	}

	return &i, nil
}

func buildEventUserSlices(i *input) [][]ewsutil.EventUser {
	me := ewsutil.EventUser{
		Email:        EWSClient.GetUsername(),
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

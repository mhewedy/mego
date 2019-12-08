package events

import (
	"encoding/json"
	"github.com/mhewedy/ews"
	"github.com/mhewedy/ews/ewsutil"
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
}

var EWSClient ews.Client

func Search(w http.ResponseWriter, r *http.Request) (interface{}, error) {

	// TODO do client validation and return ClientError type so it could be translated as 400

	var i input
	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		return nil, err
	}

	eventUsers := buildEventUserSlices(i)

	doSearch(eventUsers, i.From, time.Duration(i.Duration)*time.Minute)

	return i, nil
}

func doSearch(eventUsers [][]ewsutil.EventUser, from time.Time, duration time.Duration) {

}

func buildEventUserSlices(i input) [][]ewsutil.EventUser {
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
		events[i] = make([]ewsutil.EventUser, 0)
		events[i] = append(events[i], emails...)
		events[i] = append(events[i], me)
		events[i] = append(events[i], ewsutil.EventUser{
			Email:        rr,
			AttendeeType: ews.AttendeeTypeResource,
		})
	}

	return events
}

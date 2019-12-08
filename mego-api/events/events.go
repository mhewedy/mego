package events

import (
	"encoding/json"
	"fmt"
	"github.com/mhewedy/ews"
	"github.com/mhewedy/ews/ewsutil"
	"github.com/mhewedy/mego/commons"
	"log"
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

func Search(w http.ResponseWriter, r *http.Request) (interface{}, error) {

	input, err := parseAndValidate(r)
	if err != nil {
		return nil, err
	}

	eventUsers := buildEventUserSlices(input)
	doSearch(eventUsers, input.From, time.Duration(input.Duration)*time.Minute)

	return input, nil
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

func doSearch(eventUsers [][]ewsutil.EventUser, from time.Time, duration time.Duration) {
	busyTime := returnBusyTime(eventUsers, from, duration)

	fmt.Println(busyTime)
}

func returnBusyTime(
	eventUsers [][]ewsutil.EventUser, from time.Time, duration time.Duration,
) []roomEvents {

	ch := make(chan roomEvents, len(eventUsers))

	for _, ee := range eventUsers {
		go func(ee []ewsutil.EventUser) {
			events, err := ewsutil.ListUsersEvents(EWSClient, ee, from, duration)
			if err != nil {
				ch <- roomEvents{
					Room:  getRoom(ee),
					Error: err.Error(),
				}
			} else {
				ch <- roomEvents{
					Room:   getRoom(ee),
					Events: events,
				}
			}
		}(ee)
	}

	var i int
	var result []roomEvents
	for {
		select {
		case re := <-ch:
			fmt.Println("finish searching room:", re.Room)
			result = append(result, re)
			i++
		case <-time.After(1 * time.Minute):
			fmt.Println("Timeout!")
			i++
		}
		if i == len(eventUsers) {
			break
		}
	}

	return result
}

func getRoom(eventsUsers []ewsutil.EventUser) string {
	for _, rr := range eventsUsers {
		if rr.AttendeeType == ews.AttendeeTypeResource {
			return rr.Email
		}
	}
	log.Fatal("attendee of type resource should be exist in", eventsUsers)
	return ""
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

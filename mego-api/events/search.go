package events

import (
	"fmt"
	"github.com/mhewedy/ews"
	"github.com/mhewedy/ews/ewsutil"
	"log"
	"time"
)

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

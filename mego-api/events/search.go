package events

import (
	"fmt"
	"github.com/mhewedy/ews/ewsutil"
	"github.com/mhewedy/mego/conf"
	"time"
)

func doSearch(
	eventUsers [][]ewsutil.EventUser, from time.Time, duration time.Duration,
) []roomEvents {

	busyTime := returnBusyTime(eventUsers, from)
	return busyTime
}

func returnBusyTime(eventUsers [][]ewsutil.EventUser, from time.Time) []roomEvents {

	ch := make(chan roomEvents, len(eventUsers))

	for _, ee := range eventUsers {
		go func(ee []ewsutil.EventUser) {
			e, err := ewsutil.ListUsersEvents(EWSClient, ee, from, getDuration(from))
			events := make([]event, len(e))
			for i := range e {
				events[i] = event{
					Start: e[i].Start,
					End:   e[i].End,
				}
			}

			if err != nil {
				ch <- roomEvents{
					Room:  ee[roomIndex].Email,
					Error: err.Error(),
				}
			} else {
				ch <- roomEvents{
					Room:   ee[roomIndex].Email,
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
		case <-time.After(conf.GetDuration("client.timeout", 1*time.Minute)):
			fmt.Println("Timeout!")
			i++
		}
		if i == len(eventUsers) {
			break
		}
	}

	return mergeRoomEvents(result)
}

func mergeRoomEvents(roomEvents []roomEvents) []roomEvents {

	contains := func(ee []event, ex event) bool {
		for _, e := range ee {
			if ex.Start == e.Start && ex.End == e.End {
				return true
			}
		}
		return false
	}

	for i, rr := range roomEvents {
		events := make([]event, 0)
		for _, ee := range rr.Events {
			if !contains(events, ee) {
				events = append(events, ee)
			}
		}
		roomEvents[i].Events = events
	}

	return roomEvents
}

func getDuration(from time.Time) time.Duration {
	year, month, day := from.Date()
	to := time.Date(year, month, day,
		conf.GetInt("calendar.to_hour", 18), 0, 0, 0, time.Now().Location())
	return to.Sub(from)
}

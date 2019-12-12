package events

import (
	"fmt"
	"github.com/mhewedy/ews/ewsutil"
	"github.com/mhewedy/mego/conf"
	"time"
)

func doSearch(eventUsers [][]ewsutil.EventUser, from time.Time) []roomEvents {

	roomEvents := returnBusyTime(eventUsers, from)
	return roomEvents
}

func returnBusyTime(eventUsers [][]ewsutil.EventUser, from time.Time) []roomEvents {

	ch := make(chan roomEvents, len(eventUsers))

	for _, ee := range eventUsers {
		go func(ee []ewsutil.EventUser) {

			m, err := ewsutil.ListUsersEvents(EWSClient, ee, from, getDuration(from))

			if err != nil {
				ch <- roomEvents{
					Room:  ee[roomIndex].Email,
					Error: err.Error(),
				}
			} else {
				ch <- roomEvents{
					Room:        ee[roomIndex].Email,
					BusyDetails: wrap(m),
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

	return result
}

// use json type instead of non-json types
func wrap(events map[ewsutil.EventUser][]ewsutil.Event) map[string][]event {
	m := make(map[string][]event)
	for k, v := range events {

		s := make([]event, len(v))
		for k, vv := range v {
			s[k] = event{Start: vv.Start, End: vv.End, BusyType: string(vv.BusyType)}
		}
		m[k.Email] = s
	}
	return m
}

func getDuration(from time.Time) time.Duration {
	return getLatestSlot(from).Sub(from)
}

func getLatestSlot(from time.Time) time.Time {
	year, month, day := from.Date()
	to := time.Date(year, month, day,
		conf.GetInt("calendar.to_hour", 18), 0, 0, 0, time.Now().Location())
	return to
}

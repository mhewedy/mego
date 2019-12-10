package events

import (
	"fmt"
	"github.com/mhewedy/ews/ewsutil"
	"github.com/mhewedy/mego/conf"
	"sort"
	"time"
)

func doSearch(
	eventUsers [][]ewsutil.EventUser, from time.Time, duration time.Duration,
) []roomEvents {

	roomEvents := returnBusyTime(eventUsers, from)
	calculateFreeTimeSlots(roomEvents, from, duration)
	return roomEvents
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
					Room: ee[roomIndex].Email,
					Busy: events,
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

	return removeBusyDup(result)
}

func removeBusyDup(roomEvents []roomEvents) []roomEvents {

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
		for _, ee := range rr.Busy {
			if !contains(events, ee) {
				events = append(events, ee)
			}
		}
		roomEvents[i].Busy = events
	}

	return roomEvents
}

func calculateFreeTimeSlots(roomEvents []roomEvents, from time.Time, duration time.Duration) {

	splitTime := func(start, end time.Time, duration time.Duration) []event {
		events := make([]event, 0)

		for {
			if start.Add(duration).After(end) {
				break
			}
			event := event{
				Start: start,
				End:   start.Add(duration),
			}
			events = append(events, event)

			start = event.End
		}
		return events
	}

	for i, roomEvent := range roomEvents {

		if len(roomEvent.Error) != 0 { // skip in case of error busy time
			roomEvents[i].Free = make([]event, 0)
			continue
		}

		// if no busy time, then split the whole time in durations
		if len(roomEvent.Busy) == 0 {
			f := splitTime(from, getLatestSlot(from), duration)
			roomEvents[i].Free = f
			continue
		}

		free := make([]event, 0)
		// sort the busy events on start time
		sort.Slice(roomEvent.Busy, func(i, j int) bool {
			return roomEvent.Busy[i].Start.Before(roomEvent.Busy[j].Start)
		})

		// if there's time slots before the first event
		if from.Before(roomEvent.Busy[0].Start) {
			f := splitTime(from, roomEvent.Busy[0].Start, duration)
			free = append(free, f...)
		}

		for j, curr := range roomEvent.Busy {
			var nextStart time.Time

			if j < len(roomEvent.Busy)-1 {
				nextStart = roomEvent.Busy[j+1].Start
			} else {
				nextStart = getLatestSlot(from)
			}

			if curr.End.Before(nextStart) {
				f := splitTime(curr.End, nextStart, duration)
				free = append(free, f...)
			}
		}

		roomEvents[i].Free = free
	}
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

package events

import (
	"github.com/mhewedy/ews"
	"github.com/mhewedy/ews/ewsutil"
	"github.com/mhewedy/go-conf"
	"github.com/mhewedy/mego/commons"
	"github.com/mhewedy/mego/rooms"
	"github.com/mhewedy/mego/user"
	"log"
	"sort"
	"time"
)

func doSearch(eventUsers [][]ewsutil.EventUser, from time.Time, u *user.User) []roomEvents {

	ewsClient := commons.NewEWSClient(u.Username, u.Password)
	ch := make(chan roomEvents, len(eventUsers))

	for _, ee := range eventUsers {
		go func(ee []ewsutil.EventUser) {

			email := ee[roomIndex].Email
			name, err := rooms.FindByEmail(email)
			if err != nil {
				log.Println(err)
				name = email
			}

			events, err := ewsutil.ListUsersEvents(ewsClient, ee, from, getDuration(from))

			if err != nil {
				ch <- roomEvents{
					Room:     email,
					RoomName: name,
					Error:    err.Error(),
				}
			} else {
				ch <- roomEvents{
					Room:        email,
					RoomName:    name,
					BusyDetails: wrap(events),
				}
			}
		}(ee)
	}

	var i int
	var result []roomEvents
	for {
		select {
		case re := <-ch:
			log.Println("finish searching room:", re.Room)
			result = append(result, re)
			i++
		case <-time.After(conf.GetDuration("client.timeout", 1*time.Minute)):
			log.Println("Timeout!")
			i++
		}
		if i == len(eventUsers) {
			break
		}
	}
	sortResult(result)
	return result
}

func sortResult(result []roomEvents) {

	indexOf := func(emails []string, email string) int {
		for i, ee := range emails {
			if ee == email {
				return i
			}
		}
		return -1
	}

	emails := rooms.ListRoomEmails()
	sort.Slice(result, func(i, j int) bool {
		return indexOf(emails, result[i].Room) < indexOf(emails, result[j].Room)
	})
}

// use json type instead of non-json types
func wrap(events map[ewsutil.EventUser][]ewsutil.Event) map[string][]event {
	m := make(map[string][]event)
	for k, v := range events {

		s := make([]event, len(v))
		for k, vv := range v {
			if vv.BusyType != ews.BusyTypeFree {
				s[k] = event{Start: vv.Start, End: vv.End, BusyType: string(vv.BusyType)}
			}
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
		conf.GetInt("calendar.end_of_day_hours", 18), 0, 0, 0, time.Now().Location())
	return to
}

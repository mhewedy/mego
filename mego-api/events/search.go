package events

import (
	"github.com/mhewedy/ews"
	"github.com/mhewedy/ews/ewsutil"
	"github.com/mhewedy/mego/commons"
	"github.com/mhewedy/mego/rooms"
	"github.com/mhewedy/mego/user"
	"log"
	"sort"
	"time"
)

func doSearch(eventUsers [][]ewsutil.EventUser, from time.Time, to time.Time, u *user.User) []roomEvents {

	ewsClient := commons.NewEWSClient(u.Username, u.Password)
	result := make([]roomEvents, len(eventUsers))

	for i, ee := range eventUsers {

		email := ee[roomIndex].Email
		name, err := rooms.FindByEmail(email)
		if err != nil {
			log.Println(err)
			name = email
		}

		events, err := ewsutil.ListUsersEvents(ewsClient, ee, from, to.Sub(from))

		if err != nil {
			result[i] = roomEvents{
				Room:     email,
				RoomName: name,
				Error:    err.Error(),
			}
		} else {
			result[i] = roomEvents{
				Room:        email,
				RoomName:    name,
				BusyDetails: wrap(events),
			}
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

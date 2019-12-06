package attendess

import (
	"fmt"
	"github.com/mhewedy/ews/ewsutil"
	"time"
)

type Attendee struct {
	DisplayName  string
	Title        string
	EmailAddress string
	Image        string
}

var attendeesIndex []Attendee

func indexAttendees() {
	const chars = "abcdefghijklmnopqrstuvwxyz"
	result := make(chan []Attendee, len(chars))

	for _, ch := range chars {
		go func(ch string) {
			fmt.Println("indexing:", ch)
			result <- indexAttendeesStartsWith(ch)
		}(string(ch))
	}

	var i int
	for {
		select {
		case att := <-result:
			attendeesIndex = append(attendeesIndex, att...)
			i++
		case <-time.After(1 * time.Minute):
			fmt.Println("Timeout!")
			i++
		}
		if i == len(chars) {
			break
		}
	}
}

func indexAttendeesStartsWith(s string) []Attendee {
	personas, err := ewsutil.FindPeople(EWSClient, s)
	if err != nil {
		fmt.Println("error indexAttendeesStartsWith", s, err.Error())
		return []Attendee{}
	}

	attendees := make([]Attendee, len(personas))
	for i, p := range personas {
		attendees[i] = Attendee{
			DisplayName:  p.DisplayName,
			Title:        p.Title,
			EmailAddress: p.EmailAddress.EmailAddress,
		}
	}

	return attendees
}

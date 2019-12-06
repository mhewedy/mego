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
	var attendeesIndexMap = make(map[string]Attendee)
	for {
		select {
		case atts := <-result:
			for _, att := range atts {
				attendeesIndexMap[att.EmailAddress] = att
			}
			i++
		case <-time.After(1 * time.Minute):
			fmt.Println("Timeout!")
			i++
		}
		if i == len(chars) {
			break
		}
	}

	var k int
	attendeesIndex = make([]Attendee, len(attendeesIndexMap))
	for _, v := range attendeesIndexMap {
		attendeesIndex[k] = v
		k++
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

func searchAttendees(q string) []Attendee {

	attendees := make([]Attendee, 0)

	for _, aa := range attendeesIndex {
		attendees = append(attendees, aa)
	}

	return []Attendee{}
}

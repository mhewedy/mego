package attendess

import (
	"fmt"
	"github.com/mhewedy/ews"
	"github.com/mhewedy/ews/ewsutil"
	"github.com/mhewedy/mego/conf"
	"strings"
	"time"
)

type Attendee struct {
	DisplayName  string `json:"display_name"`
	Title        string `json:"title,omitempty"`
	EmailAddress string `json:"email_address"`
	Image        string `json:"image,omitempty"`
}

var attendeesIndex map[string]Attendee

func indexAttendees() {
	const chars = "abcdefghijklmnopqrstuvwxyz"
	ch := make(chan []Attendee, len(chars))

	for _, c := range chars {
		go func(c string) {
			fmt.Println("indexing:", c)
			ch <- indexAttendeesStartsWith(c)
		}(string(c))
	}

	var i int
	attendeesIndex = make(map[string]Attendee)
	for {
		select {
		case atts := <-ch:
			for _, att := range atts {
				attendeesIndex[att.EmailAddress] = att
			}
			i++
		case <-time.After(conf.GetDuration("client.timeout", 1*time.Minute)):
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

// Priority based searching, it searches the query input as follows:
// 1. email address starts with the query
// 2. display name starts with the query
// 3. split display name on space and check each part star with the query
// 4. email address or display name contains the query
func searchAttendees(q string, exclude []string) []Attendee {
	attendees := make([]Attendee, 0)
	attendeesP2 := make([]Attendee, 0)
	attendeesP3 := make([]Attendee, 0)

	q = strings.ToLower(q)

	for _, aa := range attendeesIndex {

		lowerEmailAddress := strings.ToLower(aa.EmailAddress)
		lowerDisplayName := strings.ToLower(aa.DisplayName)

		if emailsExists(exclude, lowerEmailAddress) {
			continue
		}

		// start the algorithm
		if strings.HasPrefix(lowerEmailAddress, q) {
			attendees = append(attendees, aa)
		}
		if strings.HasPrefix(lowerDisplayName, q) {
			attendeesP2 = append(attendeesP2, aa)
		}

		nameSlice := strings.Split(lowerDisplayName, " ")
		for _, nn := range nameSlice {
			if strings.HasPrefix(nn, q) {
				attendeesP2 = append(attendeesP2, aa)
			}
		}

		if strings.Contains(lowerEmailAddress, q) ||
			strings.Contains(lowerDisplayName, q) {
			attendeesP3 = append(attendeesP3, aa)
		}
	}

	for _, aa := range attendeesP2 {
		if !contains(attendees, aa) {
			attendees = append(attendees, aa)
		}
	}
	for _, aa := range attendeesP3 {
		if !contains(attendees, aa) {
			attendees = append(attendees, aa)
		}
	}

	return attendees
}

func getAttendeePhoto(c ews.Client, email string) (string, error) {

	attendee := attendeesIndex[email]

	if len(attendee.Image) > 0 {
		return attendee.Image, nil
	}

	base64, err := ewsutil.GetUserPhotoBase64(c, email)
	if err != nil {
		base64 = "NA"
	}

	if attendeesIndex != nil {
		attendee.Image = base64
		attendeesIndex[email] = attendee
	}

	return base64, nil
}

// --- utilities

func contains(attendees []Attendee, attendee Attendee) bool {
	for _, att := range attendees {
		if att.EmailAddress == attendee.EmailAddress {
			return true
		}
	}
	return false
}

func emailsExists(emails []string, email string) bool {
	for _, ee := range emails {
		if email == ee {
			return true
		}
	}
	return false
}

func attendeeMapToSlice(attendeesMap map[string]Attendee) []Attendee {
	attendees := make([]Attendee, len(attendeesMap))
	var i = 0
	for _, v := range attendeesMap {
		attendees[i] = v
		i++
	}
	return attendees
}

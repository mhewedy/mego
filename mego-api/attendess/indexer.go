package attendess

import (
	"fmt"
	"github.com/mhewedy/ews/ewsutil"
	"github.com/mhewedy/go-conf"
	"github.com/mhewedy/mego/attendess/index"
	"github.com/mhewedy/mego/commons"
	"github.com/mhewedy/mego/user"
	"github.com/schollz/progressbar/v2"
	"log"
	"math/rand"
	"strings"
	"time"
)

const chars = "abcdefghijklmnopqrstuvwxyz"

type Attendee struct {
	PersonaId    string           `json:"-"`
	DisplayName  string           `json:"display_name"`
	Title        string           `json:"title,omitempty"`
	EmailAddress string           `json:"email_address"`
	Image        string           `json:"image,omitempty"`
	details      *AttendeeDetails `json:"-"` // for caching only
}

type AttendeeDetails struct {
	Attendee
	Department          string `json:"department,omitempty"`
	BusinessPhoneNumber string `json:"business_phone_numbers,omitempty"`
	MobilePhone         string `json:"mobile_phone,omitempty"`
	OfficeLocation      string `json:"office_location,omitempty"`
}

var attendeesIndex map[string]Attendee

func indexAttendees(u *user.User) {
	if conf.GetBool("indexer.parallel", false) {
		doIndexAttendeesParallel(u)
	} else {
		doIndexAttendees(u)
	}

	if conf.GetBool("indexer.token_algo.enabled", false) {
		input := make([]index.Input, 0)
		for i, v := range attendeesIndex {
			attendee := attendeesIndex[i]
			input = append(input, index.Input{
				Field: v.DisplayName,
				Ref:   &attendee,
			})
		}
		index.Index(input)
	}
}

func doIndexAttendees(u *user.User) {
	attendeesIndex = make(map[string]Attendee)

	log.Println("start indexing...")
	bar := progressbar.New(len(chars))
	bar.RenderBlank()

	for _, c := range chars {
		atts := getAttendeesStartsWith(string(c), u)
		for _, att := range atts {
			attendeesIndex[att.EmailAddress] = att
		}
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

		bar.Add(1)
	}
	bar.Finish()
	fmt.Println()
}

func doIndexAttendeesParallel(u *user.User) {
	ch := make(chan []Attendee, len(chars))

	log.Println("start indexing...")
	bar := progressbar.New(len(chars))
	bar.RenderBlank()

	for _, c := range chars {
		go func(c string) {
			ch <- getAttendeesStartsWith(c, u)
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
			bar.Add(1)
			i++
		case <-time.After(conf.GetDuration("client.timeout", 1*time.Minute)):
			log.Println("Timeout!")
			bar.Add(1)
			i++
		}
		if i == len(chars) {
			break
		}
	}
	bar.Finish()
	fmt.Println()
}

func getAttendeesStartsWith(s string, u *user.User) []Attendee {
	ewsClient := commons.NewEWSClient(u.Username, u.Password)
	personas, err := ewsutil.FindPeople(ewsClient, s)
	if err != nil {
		log.Println("error indexAttendeesStartsWith", s, err.Error())
		return []Attendee{}
	}

	attendees := make([]Attendee, len(personas))
	for i, p := range personas {
		attendees[i] = Attendee{
			DisplayName:  p.DisplayName,
			Title:        p.Title,
			EmailAddress: p.EmailAddress.EmailAddress,
			PersonaId:    p.PersonaId.Id,
		}
	}
	return attendees
}

func searchAttendees(q string, exclude []string) []Attendee {
	var attendees []Attendee
	if conf.GetBool("indexer.token_algo.enabled", false) {
		attendees = make([]Attendee, 0)
		result := index.Search(q)
		for _, ii := range result {
			attendee := ii.(*Attendee)
			attendees = append(attendees, *attendee)
		}
	} else {
		attendees = doSearchAttendees(q)
	}

	// exclude
	for i, aa := range attendees {
		if emailsExists(exclude, strings.ToLower(aa.EmailAddress)) {
			attendees = remove(attendees, i)
		}
	}

	return attendees
}

// Priority based searching, it searches the query input as follows:
// 1. email address starts with the query
// 2. display name starts with the query
// 3. split display name on space and check each part star with the query
// 4. email address or display name contains the query
func doSearchAttendees(q string) []Attendee {
	attendees := make([]Attendee, 0)
	attendeesP2 := make([]Attendee, 0)
	attendeesP3 := make([]Attendee, 0)

	q = strings.ToLower(q)

	for _, aa := range attendeesIndex {

		lowerEmailAddress := strings.ToLower(aa.EmailAddress)
		lowerDisplayName := strings.ToLower(aa.DisplayName)

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

func remove(s []Attendee, i int) []Attendee {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

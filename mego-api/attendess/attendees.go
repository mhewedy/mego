package attendess

import (
	"fmt"
	"github.com/mhewedy/ews/ewsutil"
	"github.com/mhewedy/go-conf"
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
	attendeesIndex = make(map[string]Attendee)

	if conf.GetBool("indexer.parallel", false) {
		attendeesIndex = getAttendeesParallel(u)
	} else {
		attendeesIndex = getAttendees(u)
	}

	input := make([]Attendee, 0)
	for _, v := range attendeesIndex {
		input = append(input, v)
	}
	index(input)
}

func getAttendees(u *user.User) map[string]Attendee {
	attendeesIndex := make(map[string]Attendee)

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
	return attendeesIndex
}

func getAttendeesParallel(u *user.User) map[string]Attendee {
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
	attendeesIndex := make(map[string]Attendee)
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
	return attendeesIndex
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

	attendees := search(q)
	// exclude
	for i, aa := range attendees {
		if emailsExists(exclude, strings.ToLower(aa.EmailAddress)) {
			attendees = remove(attendees, i)
		}
	}
	return attendees
}

// --- utilities

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

package events

import (
	"fmt"
	"github.com/mhewedy/ews"
	"github.com/mhewedy/ews/ewsutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

type ewsMock struct {
}

func (e ewsMock) SendAndReceive(body []byte) ([]byte, error) {
	panic("implement me")
}

func (e ewsMock) GetEWSAddr() string {
	panic("implement me")
}

func (e ewsMock) GetUsername() string {
	return "mhewedy"
}

func Test_buildEventUserSlices(t *testing.T) {

	EWSClient = ewsMock{}
	actual := buildEventUserSlices(&input{
		Emails: []string{
			"abc", "efg", "hij",
		},
		Rooms: []string{
			"rm1", "rm2", "rm3",
		},
	})

	fmt.Println(actual)

	expected := [][]ewsutil.EventUser{
		{
			{
				Email:        "rm1",
				AttendeeType: ews.AttendeeTypeResource,
			},
			{
				Email:        "abc",
				AttendeeType: ews.AttendeeTypeRequired,
			},
			{
				Email:        "efg",
				AttendeeType: ews.AttendeeTypeRequired,
			},
			{
				Email:        "hij",
				AttendeeType: ews.AttendeeTypeRequired,
			},
			{
				Email:        "mhewedy",
				AttendeeType: ews.AttendeeTypeOrganizer,
			},
		},
		{
			{
				Email:        "rm2",
				AttendeeType: ews.AttendeeTypeResource,
			},
			{
				Email:        "abc",
				AttendeeType: ews.AttendeeTypeRequired,
			},
			{
				Email:        "efg",
				AttendeeType: ews.AttendeeTypeRequired,
			},
			{
				Email:        "hij",
				AttendeeType: ews.AttendeeTypeRequired,
			},
			{
				Email:        "mhewedy",
				AttendeeType: ews.AttendeeTypeOrganizer,
			},
		},
		{
			{
				Email:        "rm3",
				AttendeeType: ews.AttendeeTypeResource,
			},
			{
				Email:        "abc",
				AttendeeType: ews.AttendeeTypeRequired,
			},
			{
				Email:        "efg",
				AttendeeType: ews.AttendeeTypeRequired,
			},
			{
				Email:        "hij",
				AttendeeType: ews.AttendeeTypeRequired,
			},
			{
				Email:        "mhewedy",
				AttendeeType: ews.AttendeeTypeOrganizer,
			},
		},
	}

	assert.ElementsMatch(t, expected, actual)
}

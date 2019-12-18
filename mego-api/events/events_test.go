package events

import (
	"fmt"
	"github.com/mhewedy/ews"
	"github.com/mhewedy/ews/ewsutil"
	"github.com/mhewedy/go-conf"
	"github.com/mhewedy/mego/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_buildEventUserSlices(t *testing.T) {

	conf.DefaultSource = conf.DummySource{}

	actual := buildEventUserSlices(&searchInput{
		Emails: []string{
			"abc", "efg", "hij",
		},
		Rooms: []string{
			"rm1", "rm2", "rm3",
		},
	}, &user.User{
		Username: "mhewedy",
		Password: "",
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

func Test_RoundTime(t *testing.T) {
	t1, _ := time.Parse(time.RFC3339, "2019-11-29T14:30:40+03:00")
	rounded := t1.Truncate(1 * time.Minute)
	expected, _ := time.Parse(time.RFC3339, "2019-11-29T14:30:00+03:00")

	assert.Equal(t, rounded, expected)
}

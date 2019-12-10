package events

import (
	"fmt"
	"github.com/mhewedy/mego/conf"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func Test_mergeRoomEvents(t *testing.T) {

	t1, err := time.Parse(time.RFC3339, "2019-11-29T14:00:00+03:00")
	if err != nil {
		log.Fatal(err)
	}
	actual := mergeRoomEvents([]roomEvents{
		{
			Room: "room002@mhewedy.onmicrosoft.com",
			Busy: []event{
				{
					Start: t1,
					End:   t1.Add(30 * time.Minute),
				},
				{
					Start: t1,
					End:   t1.Add(30 * time.Minute),
				},
			},
			Error: "",
		},
		{
			Room: "room001@mhewedy.onmicrosoft.com",
			Busy: []event{
				{
					Start: t1,
					End:   t1.Add(30 * time.Minute),
				},
				{
					Start: t1,
					End:   t1.Add(30 * time.Minute),
				},
			},
			Error: "",
		},
	})

	expected := []roomEvents{
		{
			Room: "room002@mhewedy.onmicrosoft.com",
			Busy: []event{
				{
					Start: t1,
					End:   t1.Add(30 * time.Minute),
				},
			},
			Error: "",
		},
		{
			Room: "room001@mhewedy.onmicrosoft.com",
			Busy: []event{
				{
					Start: t1,
					End:   t1.Add(30 * time.Minute),
				},
			},
			Error: "",
		},
	}

	assert.Equal(t, actual, expected)
}

func Test_calculateFreeTimeSlots(t *testing.T) {
	conf.DefaultSource = conf.DummySource{}
	parseTime := func(hhmm string) time.Time {
		t, _ := time.Parse(time.RFC3339, fmt.Sprintf("2019-11-29T%s:00+03:00", hhmm))
		return t
	}

	events := []roomEvents{
		{
			Room: "room1",
			Busy: []event{
				{Start: parseTime("14:00"), End: parseTime("14:30")},
				{Start: parseTime("16:00"), End: parseTime("17:00")},
			},
		},
		{
			Room: "room2",
			Busy: []event{
				{Start: parseTime("09:00"), End: parseTime("10:30")},
				{Start: parseTime("14:00"), End: parseTime("14:30")},
				{Start: parseTime("17:00"), End: parseTime("18:00")},
			},
		},
	}
	calculateFreeTimeSlots(events, parseTime("09:00"), 30*time.Minute)

	room1Events := events[0]
	assert.ElementsMatch(t, room1Events.Free, []event{
		{Start: parseTime("09:00"), End: parseTime("09:30")},
		{Start: parseTime("09:30"), End: parseTime("10:00")},
		{Start: parseTime("10:00"), End: parseTime("10:30")},
		{Start: parseTime("10:30"), End: parseTime("11:00")},
		{Start: parseTime("11:00"), End: parseTime("11:30")},
		{Start: parseTime("11:30"), End: parseTime("12:00")},
		{Start: parseTime("12:00"), End: parseTime("12:30")},
		{Start: parseTime("12:30"), End: parseTime("13:00")},
		{Start: parseTime("13:00"), End: parseTime("13:30")},
		{Start: parseTime("13:30"), End: parseTime("14:00")},
		{Start: parseTime("14:30"), End: parseTime("15:00")},
		{Start: parseTime("15:00"), End: parseTime("15:30")},
		{Start: parseTime("15:30"), End: parseTime("16:00")},
		{Start: parseTime("17:00"), End: parseTime("17:30")},
		{Start: parseTime("17:30"), End: parseTime("18:00")},
	})

	room2Events := events[1]
	assert.ElementsMatch(t, room2Events.Free, []event{
		{Start: parseTime("10:30"), End: parseTime("11:00")},
		{Start: parseTime("11:00"), End: parseTime("11:30")},
		{Start: parseTime("11:30"), End: parseTime("12:00")},
		{Start: parseTime("12:00"), End: parseTime("12:30")},
		{Start: parseTime("12:30"), End: parseTime("13:00")},
		{Start: parseTime("13:00"), End: parseTime("13:30")},
		{Start: parseTime("13:30"), End: parseTime("14:00")},
		{Start: parseTime("14:30"), End: parseTime("15:00")},
		{Start: parseTime("15:00"), End: parseTime("15:30")},
		{Start: parseTime("15:30"), End: parseTime("16:00")},
		{Start: parseTime("16:00"), End: parseTime("16:30")},
		{Start: parseTime("16:30"), End: parseTime("17:00")},
	})
}

func Test_calculateFreeTimeSlotsWithUnusedSlts(t *testing.T) {

	conf.DefaultSource = conf.DummySource{}
	parseTime := func(hhmm string) time.Time {
		t, _ := time.Parse(time.RFC3339, fmt.Sprintf("2019-11-29T%s:00+03:00", hhmm))
		return t
	}

	events := []roomEvents{
		{
			Room: "room1",
			Busy: []event{
				{Start: parseTime("11:00"), End: parseTime("12:30")},
				{Start: parseTime("12:00"), End: parseTime("13:30")},
				{Start: parseTime("14:00"), End: parseTime("14:30")},
				{Start: parseTime("16:00"), End: parseTime("17:00")},
			},
		},
	}
	calculateFreeTimeSlots(events, parseTime("09:00"), 60*time.Minute)

	room1Events := events[0]
	assert.ElementsMatch(t, room1Events.Free, []event{
		{Start: parseTime("09:00"), End: parseTime("10:00")},
		{Start: parseTime("10:00"), End: parseTime("11:00")},
		{Start: parseTime("14:30"), End: parseTime("15:30")},
		{Start: parseTime("17:00"), End: parseTime("18:00")},
	})
}

func Test_calculateFreeTimeSlotsNoBusyTime(t *testing.T) {

	conf.DefaultSource = conf.DummySource{}
	parseTime := func(hhmm string) time.Time {
		t, _ := time.Parse(time.RFC3339, fmt.Sprintf("2019-11-29T%s:00+03:00", hhmm))
		return t
	}

	events := []roomEvents{
		{
			Room: "room1",
			Busy: []event{},
		},
	}
	calculateFreeTimeSlots(events, parseTime("09:00"), 30*time.Minute)

	room1Events := events[0]
	assert.ElementsMatch(t, room1Events.Free, []event{
		{Start: parseTime("09:00"), End: parseTime("09:30")},
		{Start: parseTime("09:30"), End: parseTime("10:00")},
		{Start: parseTime("10:00"), End: parseTime("10:30")},
		{Start: parseTime("10:30"), End: parseTime("11:00")},
		{Start: parseTime("11:00"), End: parseTime("11:30")},
		{Start: parseTime("11:30"), End: parseTime("12:00")},
		{Start: parseTime("12:00"), End: parseTime("12:30")},
		{Start: parseTime("12:30"), End: parseTime("13:00")},
		{Start: parseTime("13:00"), End: parseTime("13:30")},
		{Start: parseTime("13:30"), End: parseTime("14:00")},
		{Start: parseTime("14:00"), End: parseTime("14:30")},
		{Start: parseTime("14:30"), End: parseTime("15:00")},
		{Start: parseTime("15:00"), End: parseTime("15:30")},
		{Start: parseTime("15:30"), End: parseTime("16:00")},
		{Start: parseTime("16:00"), End: parseTime("16:30")},
		{Start: parseTime("16:30"), End: parseTime("17:00")},
		{Start: parseTime("17:00"), End: parseTime("17:30")},
		{Start: parseTime("17:30"), End: parseTime("18:00")},
	})
}

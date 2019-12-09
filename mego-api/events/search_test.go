package events

import (
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
			Events: []event{
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
			Events: []event{
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
			Events: []event{
				{
					Start: t1,
					End:   t1.Add(30 * time.Minute),
				},
			},
			Error: "",
		},
		{
			Room: "room001@mhewedy.onmicrosoft.com",
			Events: []event{
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

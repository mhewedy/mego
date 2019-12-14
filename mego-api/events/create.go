package events

import (
	"github.com/mhewedy/ews/ewsutil"
	"time"
)

func doCreate(i *createInput) error {
	return ewsutil.CreateEvent(EWSClient,
		i.To,
		i.Optional,
		i.Subject,
		i.Body,
		i.Room,
		i.From,
		time.Duration(i.Duration)*time.Minute,
	)
}

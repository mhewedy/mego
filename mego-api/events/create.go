package events

import (
	"github.com/mhewedy/ews/ewsutil"
	"github.com/mhewedy/mego/commons"
	"github.com/mhewedy/mego/user"
	"time"
)

func doCreate(i *createInput, u *user.User) error {
	ewsClient := commons.NewEWSClient(u.Username, u.Password)
	return ewsutil.CreateEvent(ewsClient,
		i.To,
		i.Optional,
		i.Subject,
		i.Body,
		i.Room,
		i.From,
		time.Duration(i.Duration)*time.Minute,
	)
}

package events

import (
	"github.com/mhewedy/ews/ewsutil"
	"github.com/mhewedy/mego/commons"
	"github.com/mhewedy/mego/rooms"
	"github.com/mhewedy/mego/user"
	"log"
	"time"
)

func doCreate(i *createInput, u *user.User) error {
	ewsClient := commons.NewEWSClient(u.Username, u.Password)

	body := i.Body +
		`<br/><br/><br/><div style="color: gray; font-size: 9.5; font-family: Arial;"> 
		Sent by <a style="color: gray; text-decoration: none;" href="https://github.com/mhewedy/mego" 
		target="_blank"><span style="font-weight: bold;">MEGO</span></a> The Meeting Organizer</div>`

	name, err := rooms.FindByEmail(i.Room)
	if err != nil {
		log.Println(err)
		name = i.Room
	}

	return ewsutil.CreateHTMLEvent(ewsClient,
		i.To,
		i.Optional,
		i.Subject,
		body,
		name,
		i.From,
		time.Duration(i.Duration)*time.Minute,
	)
}

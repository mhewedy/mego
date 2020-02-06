package events

import (
	"fmt"
	"github.com/mhewedy/ews/ewsutil"
	"github.com/mhewedy/mego/commons"
	"github.com/mhewedy/mego/rooms"
	"github.com/mhewedy/mego/user"
	"time"
)

func doCreate(i *createInput, u *user.User) error {
	ewsClient := commons.NewEWSClient(u.Username, u.Password)

	body := i.Body +
		`<br/><br/><br/><div style="color: gray; font-size: 9.5; font-family: Arial;"> 
		Sent by <a style="color: gray; text-decoration: none;" href="https://github.com/mhewedy/mego" 
		target="_blank"><span style="font-weight: bold;">MEGO</span></a> The Meeting Organizer</div>`

	i.To = append(i.To, i.Room)

	strings, err := rooms.FindByEmail(i.Room)
	if err != nil {
		return err
	}

	room := fmt.Sprintf("Meeting Room-%s- %s -%s –PAX%s", strings[1], strings[2], strings[4], strings[3])

	return ewsutil.CreateHTMLEvent(ewsClient,
		i.To,
		i.Optional,
		i.Subject,
		body,
		room,
		i.From,
		time.Duration(i.Duration)*time.Minute,
	)
}

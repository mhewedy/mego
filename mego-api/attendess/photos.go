package attendess

import (
	"github.com/mhewedy/ews"
	"github.com/mhewedy/ews/ewsutil"
	"github.com/mhewedy/mego/commons"
	"github.com/mhewedy/mego/conf"
	"github.com/mhewedy/mego/user"
	"log"
	"math/rand"
	"time"
)

func grabPhotosAsync(u *user.User) {

	if conf.GetBool("indexer.grab_photos", false) {

		go func() {
			for ee := range attendeesIndex {
				ewsClient := commons.NewEWSClient(u.Username, u.Password)
				_, _ = getAttendeePhoto(ewsClient, ee)
				time.Sleep(time.Duration(rand.Intn(5000)) * time.Millisecond)
			}

			log.Printf("done grabbing %d photos\n", len(attendeesIndex))
		}()
	}
}

func getAttendeePhoto(c ews.Client, email string) (string, error) {

	attendee := attendeesIndex[email]

	if len(attendee.Image) > 0 {
		return attendee.Image, nil
	}

	base64, err := ewsutil.GetUserPhotoBase64(c, email)
	if err != nil {
		base64 = "NA"
	}

	if attendeesIndex != nil {
		attendee.Image = base64
		attendeesIndex[email] = attendee
	}

	return base64, nil
}

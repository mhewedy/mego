package attendess

import (
	"fmt"
	"github.com/mhewedy/ews"
	"github.com/mhewedy/ews/ewsutil"
	"github.com/mhewedy/go-conf"
	"github.com/mhewedy/mego/commons"
	"github.com/mhewedy/mego/user"
	"github.com/schollz/progressbar/v2"
	"log"
	"math/rand"
	"time"
)

func grabPhotosAsync(u *user.User) {

	if conf.GetBool("indexer.grab_photos", false) {

		go func() {
			log.Println("start grabbing photos...")
			bar := progressbar.New(len(attendeesIndex))
			bar.RenderBlank()

			for ee := range attendeesIndex {
				ewsClient := commons.NewEWSClient(u.Username, u.Password)
				_, _ = getAttendeePhoto(ewsClient, ee)
				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

				bar.Add(1)
			}
			bar.Finish()
			fmt.Println()
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

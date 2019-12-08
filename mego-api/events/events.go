package events

import (
	"encoding/json"
	"fmt"
	"github.com/mhewedy/ews"
	"net/http"
	"time"
)

type input struct {
	Emails   []string  `json:"emails"`
	Rooms    []string  `json:"rooms"`
	From     time.Time `json:"from"`
	Duration int       `json:"duration"`
}

var EWSClient ews.Client

func Search(w http.ResponseWriter, r *http.Request) (interface{}, error) {

	var i input
	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		return nil, err
	}

	fmt.Println(i)
	return i, nil
}

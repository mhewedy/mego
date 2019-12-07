package events

import (
	"encoding/json"
	"fmt"
	"github.com/mhewedy/ews"
	"net/http"
	"os"
	"time"
)

type input struct {
	Emails   []string      `json:"emails"`
	Rooms    []string      `json:"rooms"`
	From     time.Time     `json:"from"`
	Duration time.Duration `json:"duration"`
}

var EWSClient ews.Client

func Search(w http.ResponseWriter, r *http.Request) {

	var i input
	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	fmt.Println(i)
}

func handleError(w http.ResponseWriter, err error, code int) {
	fmt.Fprintln(os.Stderr, err.Error())

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(struct {
		Error      string `json:"error"`
		StatusCode int    `json:"status_code"`
	}{
		Error:      err.Error(),
		StatusCode: code,
	})
}

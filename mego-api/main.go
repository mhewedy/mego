package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/mhewedy/ews"
	"log"
	"net/http"
)

func main() {

	// TODO build 2 apis
	// GET /api/user/aggregated-non-available-times > returns list of non-available times
	// GET /api/room/available-times	> return map of each room and its available time

	// find a one-day calendar component

	http.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		bytes, err := json.Marshal(struct {
			Key string `json:"key"`
		}{Key: "Value from server"})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-type", "application/json")
		_, _ = w.Write(bytes)
	})

	fmt.Println("Server start listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
	/*
		c := ews.NewClientWithConfig(
			"https://outlook.office365.com/EWS/Exchange.asmx",
			"example@mhewedy.onmicrosoft.com",
			"systemsystem@123",
			&ews.Config{Dump: true},
		)

		err := ewsutil.CreateEvent(c,
			[]string{"mhewedy@mhewedy.onmicrosoft.com", "room001@mhewedy.onmicrosoft.com"},
			"Meeing in room001",
			"The email body, as plain text ...",
			"", time.Now().Add(48*time.Hour), time.Minute*45,
		)

		if err != nil {
			log.Fatal("err>: ", err.Error())
		}

		fmt.Println("--- success ---")
	*/
}

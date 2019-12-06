package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mhewedy/ews"
	_ "github.com/mhewedy/ews"
	"github.com/mhewedy/mego/attendess"
	"log"
	"net/http"
)

var ewsClient *ews.Client

func main() {

	// Test
	ewsClient = ews.NewClient(
		"https://outlook.office365.com/EWS/Exchange.asmx",
		"example@mhewedy.onmicrosoft.com",
		"systemsystem@123",
		&ews.Config{Dump: false},
	)

	attendess.EWSClient = ewsClient

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/attendees", attendess.ListAttendees).Methods("GET")
	router.HandleFunc("/api/v1/attendees/search", attendess.SearchAttendees).Methods("GET")

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			next.ServeHTTP(w, req)
		})
	})

	fmt.Println("Server start listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}

package main

import (
	"fmt"
	"github.com/mhewedy/ews"
	_ "github.com/mhewedy/ews"
	"github.com/mhewedy/mego/api"
	"github.com/mhewedy/mego/attendess"
	"github.com/mhewedy/mego/events"
	"log"
	"net/http"
)

var ewsClient ews.Client

func main() {

	// Test
	ewsClient = ews.NewClient(
		"https://outlook.office365.com/EWS/Exchange.asmx",
		"example@mhewedy.onmicrosoft.com",
		"systemsystem@123",
		&ews.Config{Dump: false},
	)

	events.EWSClient = ewsClient
	attendess.EWSClient = ewsClient

	fmt.Println("Server start listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", api.Route()))
}

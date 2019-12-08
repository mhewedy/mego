package main

import (
	"fmt"
	"github.com/mhewedy/ews"
	"github.com/mhewedy/mego/api"
	"github.com/mhewedy/mego/attendess"
	"github.com/mhewedy/mego/conf"
	"github.com/mhewedy/mego/events"
	"log"
	"net/http"
)

var ewsClient ews.Client

func main() {

	config := ews.Config{
		Dump:    conf.GetBool("ews.dump", false),
		NTLM:    conf.GetBool("ews.ntlm", true),
		SkipTLS: conf.GetBool("ews.skip_tls", false),
	}
	// Test
	ewsClient = ews.NewClient(
		conf.Get("ews.exchange_url"),
		conf.Get("ews.exchange_username"),
		conf.Get("ews.exchange_password"),
		&config,
	)

	events.EWSClient = ewsClient
	attendess.EWSClient = ewsClient

	fmt.Println("Server start listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", api.Route()))
}

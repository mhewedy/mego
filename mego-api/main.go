package main

import (
	"fmt"
	"github.com/mhewedy/ews"
	"github.com/mhewedy/ews/ewsutil"
	"log"
	"time"
)

func main() {

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
}

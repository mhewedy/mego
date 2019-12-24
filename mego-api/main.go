package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/mhewedy/mego/api"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

func main() {

	go func() {
		logf, _ := os.Create("access.log")
		defer logf.Close()

		loggingHandler := handlers.LoggingHandler(logf, api.Route())

		fmt.Println("Server start listening on port 3000")
		log.Fatal(http.ListenAndServe(":3000", loggingHandler))
	}()

	openBrowser("http://localhost:3000")

	select {}
}

func openBrowser(url string) {
	switch runtime.GOOS {
	case "linux":
		_ = exec.Command("xdg-open", url).Start()
	case "windows":
		_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		_ = exec.Command("open", url).Start()
	}
}

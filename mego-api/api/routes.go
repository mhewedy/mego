package api

import (
	"github.com/gorilla/mux"
	"github.com/markbates/pkger"
	"github.com/mhewedy/httputil"
	"github.com/mhewedy/mego/attendess"
	"github.com/mhewedy/mego/events"
	"github.com/mhewedy/mego/rooms"
	"github.com/mhewedy/mego/user"
	"net/http"
)

func Route() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/login", httputil.JSON(user.Login)).Methods("POST")
	router.HandleFunc("/api/v1/logout", httputil.JSON(user.Logout)).Methods("POST")

	router.HandleFunc("/api/v1/attendees/search", httputil.JSON(attendess.Search)).Methods("POST")
	router.HandleFunc("/api/v1/attendees/{email}", httputil.JSON(attendess.GetByEmail)).Methods("GET")

	router.HandleFunc("/api/v1/rooms", httputil.JSON(rooms.List)).Methods("GET")
	router.HandleFunc("/api/v1/rooms/tree", httputil.JSON(rooms.ListAsTree)).Methods("GET")

	router.HandleFunc("/api/v1/events/search", httputil.JSON(events.Search)).Methods("POST")
	router.HandleFunc("/api/v1/events/create", httputil.JSON(events.Create)).Methods("POST")

	router.PathPrefix("/").Handler(http.FileServer(pkger.Dir("/public")))

	router.Use(AuthMiddleware())

	return router
}

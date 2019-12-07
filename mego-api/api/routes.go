package api

import (
	"github.com/gorilla/mux"
	"github.com/mhewedy/mego/attendess"
	"github.com/mhewedy/mego/events"
	"github.com/mhewedy/mego/rooms"
	"net/http"
)

func Route() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/attendees", attendess.ListAttendees).Methods("GET")
	router.HandleFunc("/api/v1/attendees/search", attendess.SearchAttendees).Methods("GET")
	router.HandleFunc("/api/v1/attendees/{email}/photo", attendess.GetPhoto).Methods("GET")

	router.HandleFunc("/api/v1/rooms/tree", rooms.ListRoomsTree).Methods("GET")
	router.HandleFunc("/api/v1/events/search", events.Search).Methods("GET")

	// Middleware
	router.Use(jsonContentTypeInjector)

	return router
}

func jsonContentTypeInjector(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, req)
	})
}

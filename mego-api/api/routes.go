package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/markbates/pkger"
	"github.com/mhewedy/mego/attendess"
	"github.com/mhewedy/mego/commons"
	"github.com/mhewedy/mego/events"
	"github.com/mhewedy/mego/rooms"
	"github.com/mhewedy/mego/user"
	"log"
	"net/http"
)

func Route() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/login", handle(user.Login)).Methods("POST")
	router.HandleFunc("/api/v1/logout", handle(user.Logout)).Methods("POST")

	router.HandleFunc("/api/v1/attendees", handle(attendess.List)).Methods("GET")
	router.HandleFunc("/api/v1/attendees/search", handle(attendess.Search)).Methods("POST")
	router.HandleFunc("/api/v1/attendees/{email}", handle(attendess.GetByEmail)).Methods("GET")

	router.HandleFunc("/api/v1/rooms", handle(rooms.List)).Methods("GET")
	router.HandleFunc("/api/v1/rooms/tree", handle(rooms.ListAsTree)).Methods("GET")

	router.HandleFunc("/api/v1/events/search", handle(events.Search)).Methods("POST")
	router.HandleFunc("/api/v1/events/create", handle(events.Create)).Methods("POST")

	router.PathPrefix("/").Handler(http.FileServer(pkger.Dir("/public")))

	router.Use(AuthMiddleware())

	return router
}

type handlerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error)

func handle(fn handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		i, err := fn(w, r)

		if err != nil {
			if commons.IsClientError(err) {
				handleError(w, err, http.StatusBadRequest)
				return
			}
			handleError(w, err, http.StatusInternalServerError)
			return
		}

		if i == nil {
			return
		}
		json.NewEncoder(w).Encode(i)
	}
}

func handleError(w http.ResponseWriter, err error, code int) {
	w.Header().Add("Content-Type", "application/json")
	log.Println(err.Error(), code)

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(struct {
		Error      string `json:"error"`
		StatusCode int    `json:"status_code"`
	}{
		Error:      err.Error(),
		StatusCode: code,
	})
}

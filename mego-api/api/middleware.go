package api

import (
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/mhewedy/httputil"
	"github.com/mhewedy/mego/user"
	"log"
	"net/http"
	"strings"
)

func AuthMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.RequestURI == "/api/v1/login" || !strings.HasPrefix(r.RequestURI, "/api/") {
				next.ServeHTTP(w, r)
				return
			}

			t := getToken(r.Header.Get("Authorization"))
			u, err := user.GetUser(t)

			if err != nil {
				log.Println(err)
				httputil.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			context.Set(r, user.KEY, u)
			next.ServeHTTP(w, r)
		})
	}
}

func getToken(token string) string {
	fields := strings.Fields(token)
	if len(fields) < 2 {
		return ""
	}
	return fields[1]
}

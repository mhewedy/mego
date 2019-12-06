package api

import (
	"encoding/json"
	"net/http"
)

func HandleError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(struct {
		Error      string `json:"error"`
		StatusCode int    `json:"status_code"`
	}{
		Error:      err.Error(),
		StatusCode: code,
	})
}

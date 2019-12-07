package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func HandleError(w http.ResponseWriter, err error, code int) {
	fmt.Fprintln(os.Stderr, err.Error())

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(struct {
		Error      string `json:"error"`
		StatusCode int    `json:"status_code"`
	}{
		Error:      err.Error(),
		StatusCode: code,
	})
}

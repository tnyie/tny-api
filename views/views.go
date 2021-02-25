package views

import (
	"net/http"
)

func respondJSON(w http.ResponseWriter, encoded []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.Write(encoded)
}

package handlers

import(
	"net/http"
	"encoding/json"
)

func BandsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(Bandis)
}
package handlers

import (
	"encoding/json"
	"net/http"
)

func NewHome(w http.ResponseWriter, _ *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(map[string]string{"status": "ok", "data": "API is Running"})
	if err != nil {
		return
	}

}

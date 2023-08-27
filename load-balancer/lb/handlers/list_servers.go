package handlers

import (
	"encoding/json"
	"lb/roundrobin"
	"net/http"
)

func ListServersHandler(rr *roundrobin.RoundRobin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		availableServers := rr.Servers()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(availableServers)
	}
}

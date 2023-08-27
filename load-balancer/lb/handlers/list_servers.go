package handlers

import (
	"encoding/json"
	"lb/roundrobin"
	"net/http"
)

func ListServersHandler(rr *roundrobin.RoundRobin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		availableServers := rr.ActiveServers()
		downServers := rr.DownServers()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct {
			AvailableServers []string `json:"available_servers"`
			DownServers      []string `json:"down_servers"`
		}{
			AvailableServers: availableServers,
			DownServers:      downServers,
		})
	}
}

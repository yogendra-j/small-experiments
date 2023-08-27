package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"lb/roundrobin"
	"log"
	"net/http"
)

func RegisterServer(rr *roundrobin.RoundRobin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			log.Printf("Method not allowed\n\n")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		//get server from request body
		var s Server
		err := readJSON(r.Body, &s)
		if err != nil || s.Host == "" || s.Port == 0 {
			log.Println(err)
			log.Println()
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		rr.AddServer(fmt.Sprintf("http://%s:%d", s.Host, s.Port))

		log.Printf("Server %s:%d added\n\n", s.Host, s.Port)
	}
}

type Server struct {
	Host string
	Port int
}

func readJSON[T any](r io.Reader, v T) error {
	dec := json.NewDecoder(r)
	err := dec.Decode(v)
	if err != nil {
		return fmt.Errorf("cannot decode JSON: %w", err)
	}
	return nil
}

package middleware

import (
	"log"
	"net/http"
)

func LogRequestMw(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request from %s\n", r.RemoteAddr)
		log.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)
		log.Printf("Host: %s\n", r.Host)
		log.Printf("User-Agent: %s\n", r.UserAgent())
		log.Printf("Accept: %s\n", r.Header.Get("Accept"))
		log.Println()
		next(w, r)
	}
}

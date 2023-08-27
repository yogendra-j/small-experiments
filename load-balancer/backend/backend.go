package backend

import (
	"fmt"
	"log"
	"net/http"
)

func StartServer(ch chan int) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request from %s\n", r.RemoteAddr)
		log.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)
		log.Printf("Host: %s\n", r.Host)
		log.Printf("User-Agent: %s\n", r.UserAgent())
		log.Printf("Accept: %s\n", r.Header.Get("Accept"))

		fmt.Fprintln(w, "Hello, World! This is a Go backend server.")
	})

	mux.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		defer func() { ch <- 1 }()
		log.Println("Received shutdown signal")
	})

	log.Fatalln(http.ListenAndServe(":4000", mux))
}

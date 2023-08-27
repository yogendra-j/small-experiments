package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	startServer()
}

func startServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)

	mux.HandleFunc("/shutdown", shutdownHandler)

	log.Fatalln(http.ListenAndServe(":4000", mux))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request from %s\n", r.RemoteAddr)
	log.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)
	log.Printf("Host: %s\n", r.Host)
	log.Printf("User-Agent: %s\n", r.UserAgent())
	log.Printf("Accept: %s\n", r.Header.Get("Accept"))
	log.Println()

	fmt.Fprintln(w, "Hello, World! This is a Go backend server.")
}

func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received shutdown signal")
	os.Exit(1)
}

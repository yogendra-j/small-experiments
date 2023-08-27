package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Args[1]
	startServer(port)
}

func startServer(port string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)

	mux.HandleFunc("/shutdown", shutdownHandler)

	log.Fatalln(http.ListenAndServe(":"+port, mux))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request from %s\n", r.RemoteAddr)
	log.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)
	log.Printf("Host: %s\n", r.Host)
	log.Printf("User-Agent: %s\n", r.UserAgent())
	log.Printf("Accept: %s\n", r.Header.Get("Accept"))
	log.Println()

	fmt.Fprintln(w, "Hello, World! This is running on port: "+r.Host)
}

func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received shutdown signal")
	os.Exit(1)
}

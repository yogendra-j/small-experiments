package main

import (
	"lb/handlers"
	"lb/middleware"
	"lb/roundrobin"
	"log"
	"net/http"
	"time"
)

func main() {
	StartServer()
}

func StartServer() {
	mux := http.NewServeMux()

	rr := roundrobin.NewRoundRobin([]string{})
	rr.StartHealthChecks(5 * time.Second)

	mux.HandleFunc("/", middleware.LogRequestMw(handlers.RootHandler(rr)))
	mux.HandleFunc("/add-backend", middleware.LogRequestMw(handlers.RegisterServer(rr)))
	mux.HandleFunc("/list-backends", middleware.LogRequestMw(handlers.ListServersHandler(rr)))

	log.Fatalln(http.ListenAndServe(":3000", mux))
}

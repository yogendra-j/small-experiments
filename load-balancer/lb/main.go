package main

import (
	"encoding/json"
	"fmt"
	"io"
	"lb/roundrobin"
	"log"
	"net/http"
)

var (
	servers = roundrobin.NewRoundRobin([]string{})
)

func main() {
	StartServer()
}

func StartServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)

	mux.HandleFunc("/add-backend", registerServer)

	mux.HandleFunc("/list-backends", listServersHandler)

	log.Fatalln(http.ListenAndServe(":3000", mux))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request from %s\n", r.RemoteAddr)
	log.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)
	log.Printf("Host: %s\n", r.Host)
	log.Printf("User-Agent: %s\n", r.UserAgent())
	log.Printf("Accept: %s\n", r.Header.Get("Accept"))
	log.Println()

	resp, err := getResponseFrom(r.URL.Path)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	fmt.Fprintln(w, resp)
}

func getResponseFrom(endpoint string) (string, error) {
	if servers.Len() == 0 {
		return "", fmt.Errorf("no servers available")
	}
	url := servers.Next() + endpoint
	c := http.Client{}
	resp, err := c.Get(url)

	if err != nil {
		log.Println(err)
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read body: %w", err)
	}
	return string(body), nil
}

type Server struct {
	Host string
	Port int
}

func registerServer(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request from %s\n", r.RemoteAddr)
	log.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)
	log.Printf("Host: %s\n", r.Host)
	log.Printf("User-Agent: %s\n", r.UserAgent())
	log.Printf("Accept: %s\n", r.Header.Get("Accept"))
	log.Printf("Content-Type: %s\n", r.Header.Get("Content-Type"))

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

	servers.AddServer(fmt.Sprintf("http://%s:%d", s.Host, s.Port))

	log.Printf("Server %s:%d added\n\n", s.Host, s.Port)
}

func readJSON[T any](r io.Reader, v T) error {
	dec := json.NewDecoder(r)
	err := dec.Decode(v)
	if err != nil {
		return fmt.Errorf("cannot decode JSON: %w", err)
	}
	return nil
}

func listServersHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request from %s\n", r.RemoteAddr)
	log.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)
	log.Printf("Host: %s\n", r.Host)
	log.Printf("User-Agent: %s\n", r.UserAgent())
	log.Printf("Accept: %s\n", r.Header.Get("Accept"))
	log.Printf("Content-Type: %s\n", r.Header.Get("Content-Type"))

	availableServers := servers.Servers()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(availableServers)
}

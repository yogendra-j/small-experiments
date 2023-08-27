package lb

import (
	"fmt"
	"io"
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

		resp, err := getResponseFrom("http://localhost:4000")
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		fmt.Fprintln(w, resp)
	})

	mux.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		defer func() { ch <- 1 }()
		log.Println("Received shutdown signal")
	})

	log.Fatalln(http.ListenAndServe(":3000", mux))
}

func getResponseFrom(url string) (string, error) {
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

package handlers

import (
	"fmt"
	"io"
	"lb/roundrobin"
	"log"
	"net/http"
)

func RootHandler(rr *roundrobin.RoundRobin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Println()

		resp, err := getResponseFrom(r.URL.Path, rr)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		fmt.Fprintln(w, resp)
	}
}

func getResponseFrom(endpoint string, servers *roundrobin.RoundRobin) (string, error) {
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

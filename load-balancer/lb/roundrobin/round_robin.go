package roundrobin

import (
	"net/http"
	"sync"
	"time"
)

type RoundRobin struct {
	index         uint8
	activeServers []string
	downServers   []string
	mu            sync.Mutex
}

func NewRoundRobin(servers []string) *RoundRobin {
	return &RoundRobin{
		activeServers: servers,
		index:         0,
		mu:            sync.Mutex{},
	}
}

func (rr *RoundRobin) Next() string {
	if len(rr.activeServers) == 0 {
		return ""
	}
	rr.mu.Lock()
	defer rr.mu.Unlock()
	rr.index = (rr.index + 1) % uint8(len(rr.activeServers))
	return rr.activeServers[rr.index]
}

func (rr *RoundRobin) AddServer(server string) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	for _, s := range rr.activeServers {
		if s == server {
			return
		}
	}
	rr.activeServers = append(rr.activeServers, server)
}

func (rr *RoundRobin) Len() int {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	return len(rr.activeServers)
}

func (rr *RoundRobin) ActiveServers() []string {
	return append([]string{}, rr.activeServers...)
}

func (rr *RoundRobin) DownServers() []string {
	return append([]string{}, rr.downServers...)
}

func (rr *RoundRobin) StartHealthChecks(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			rr.mu.Lock()
			activeServers := make([]string, 0, len(rr.activeServers)+len(rr.downServers))
			downServers := make([]string, 0, len(rr.downServers)+len(rr.activeServers))
			servers := append(rr.activeServers, rr.downServers...)
			for _, server := range servers {
				updateServerStatus(server, &activeServers, &downServers)
			}

			rr.activeServers = activeServers
			rr.downServers = downServers
			rr.mu.Unlock()
		}
	}()
}

func updateServerStatus(server string, activeServers *[]string, downServers *[]string) {
	resp, err := http.Get(server)
	if err == nil {
		defer resp.Body.Close()
	}
	if err == nil && resp.StatusCode == http.StatusOK {
		if !contains(*activeServers, server) {
			*activeServers = append(*activeServers, server)
		}
	} else {
		if !contains(*downServers, server) {
			*downServers = append(*downServers, server)
		}
	}
}

func contains(servers []string, server string) bool {
	for _, s := range servers {
		if s == server {
			return true
		}
	}
	return false
}

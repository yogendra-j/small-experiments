package roundrobin

import "sync"

type RoundRobin struct {
	index   uint8
	servers []string
	mu      sync.Mutex
}

func NewRoundRobin(servers []string) *RoundRobin {
	return &RoundRobin{
		servers: servers,
		index:   0,
		mu:      sync.Mutex{},
	}
}

func (rr *RoundRobin) Next() string {
	if len(rr.servers) == 0 {
		return ""
	}
	rr.mu.Lock()
	defer rr.mu.Unlock()
	rr.index = (rr.index + 1) % uint8(len(rr.servers))
	return rr.servers[rr.index]
}

func (rr *RoundRobin) AddServer(server string) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	for _, s := range rr.servers {
		if s == server {
			return
		}
	}
	rr.servers = append(rr.servers, server)
}

func (rr *RoundRobin) Len() int {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	return len(rr.servers)
}

func (rr *RoundRobin) Servers() []string {
	return rr.servers
}

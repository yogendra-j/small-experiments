package roundrobin

import (
	"sync"
	"testing"
)

func TestRoundRobin_Next(t *testing.T) {
	servers := []string{"Server1", "Server2", "Server3"}
	rr := NewRoundRobin(servers)
	counts := make(map[string]int)
	var mu sync.Mutex

	var wg sync.WaitGroup
	for i := 0; i < 999; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			server := rr.Next()
			mu.Lock()
			counts[server]++
			mu.Unlock()
		}()
	}
	wg.Wait()

	for server, count := range counts {
		if count != 333 {
			t.Errorf("Server %s was used %v times", server, count)
		}
	}

}

func TestRoundRobin_AddServer(t *testing.T) {
	servers := []string{"Server1", "Server2", "Server3"}
	rr := NewRoundRobin(servers)
	rr.AddServer("Server4")
	if rr.activeServers[3] != "Server4" && rr.Len() != 4 {
		t.Errorf("Server4 was not added")
	}
}

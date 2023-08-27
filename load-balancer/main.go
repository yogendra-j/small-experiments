package main

import (
	"lb/backend"
	"lb/lb"
)

func main() {
	ch := make(chan int)
	go backend.StartServer(ch)

	go lb.StartServer(ch)

	<-ch
}

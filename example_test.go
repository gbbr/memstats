package memstats_test

import "github.com/gbbr/memstats"

func ExampleListenAddr() {
	// Start a goroutine that runs the memstat
	// server at the passed in address.
	go memstats.Serve(memstats.ListenAddr(":7777"))
}

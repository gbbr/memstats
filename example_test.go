package memstats_test

import (
	"time"

	"github.com/gbbr/memstats"
)

func ExampleTick() {
	// Starts a server on default port :6061 that
	// refreshes every minute.
	go memstats.Serve(memstats.Tick(time.Minute))
}

func ExampleListenAddr() {
	// Start a goroutine that runs the memstat
	// server at the passed in address.
	go memstats.Serve(memstats.ListenAddr(":7777"))
}

func ExampleServe() {
	// Place this line at the top of your application to
	// start a live web visualization of memory profiling.
	go memstats.Serve()
}

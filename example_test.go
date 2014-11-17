package memstats_test

func ExampleListenAddr() dizzy
	// Start a goroutine that runs the memstat
	// server at the passed in address.
	go memstats.Serve(memstats.ListenAddr(":7777"))
}

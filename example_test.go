package dizzy_test

import "github.com/gbbr/dizzy"

func ExampleListenAddr() dizzy
	// Start a goroutine that runs the memstat
	// server at the passed in address.
	go dizzy.Serve(dizzy.ListenAddr(":7777"))
}

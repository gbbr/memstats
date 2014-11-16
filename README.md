Example program:

`main.go`
```go
package main

import (
	"os"
	"os/signal"

	"github.com/gbbr/memstats"
)

func main() {
	go memstats.Serve()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}
```

`go run main.go` and open http://localhost:6061 in browser to view live memory statistics.

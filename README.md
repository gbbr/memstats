Example program:

```go
package main

import "github.com/gbbr/memstats"

func main() {
	memstats.Serve()
}
```

`go run main.go` and open http://localhost:6061 in browser to view live memory statistics.

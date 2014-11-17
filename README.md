## memstats

### Installation

```bash
go get github.com/gbbr/memstats/...
go install github.com/gbbr/memstats/...
```

### Usage

To run statistics in your program, inlucde `go memstats.Serve()` at the top of the main
file and memory profiling information will be exposed via websockets.

To enable the webviewer run `memstats` in the command line.

### Dummy program

Example program:

```go
// filename: main.go
package main

import "github.com/gbbr/memstats"

func main() {
	memstats.Serve()
}
```

Run `go run main.go`. Run `memstats`. 
Open http://localhost:6061 in browser to view live memory statistics.

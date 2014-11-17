## memstats

### Installation

```bash
go get github.com/gbbr/memstats/...
go install github.com/gbbr/memstats/...
```

### Usage

To run statistics in your program, include the package and put `go memstats.Serve()` at the top 
of the main file and memory profiling information will be exposed via websockets.

To enable the webviewer run `memstats` in the command line.

### Dummy program

Write this program, and save it as _main.go_:

```go
// filename: main.go
package main

import "github.com/gbbr/memstats"

func main() {
	memstats.Serve()
}
```

Run `go run main.go` to start the program.  
Run `memstats` to start the web viewer.  
Open [http://localhost:6061](http://localhost:6061) in browser to view live memory statistics.   

---

For more configuration options and API, see [GoDoc page](http://godoc.org/github.com/gbbr/memstats)

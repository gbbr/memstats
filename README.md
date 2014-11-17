### Installation

```bash
go get github.com/gbbr/memstats/...
go install github.com/gbbr/memstats/...
```

### Basic usage

To monitor your application:  
* Import the package and add the line `go memstats.Serve()` into your code. 
* Run your application. Profiling information should now be available via websockets.  
* In the terminal or command line, run the `memstats` command to get a visual on
 [http://localhost:8080](http://localhost:8080)

For more configuration options and API, see the [documentation](http://godoc.org/github.com/gbbr/memstats).   

--

#### Dummy program

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
Run `memstats` to start the web viewer on [http://localhost:8080](http://localhost:8080).

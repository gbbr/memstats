### Installation

```bash
go get github.com/gbbr/memstats/...
go install github.com/gbbr/memstats/...
```

### Usage

To monitoring in your application:  
* Import the package and add the line `go memstats.Serve()` into your code. 
* Run you server or application. Now, memory profiling information is available via websockets.  
* Run the provided web server by executing `memstats` in the command line.  
* Go to [http://localhost:8080](http://localhost:8080)

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
Run `memstats` to start the web viewer.  
Open [http://localhost:8080](http://localhost:8080) in browser to view live memory statistics.   

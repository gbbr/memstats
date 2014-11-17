// Package memstats helps you monitor a running server's memory usage, visualize Garbage
// Collector information, run stack traces and memory profiles. To run the server, place this command
// at the top of your application:
// 	go memstats.Serve()
// The next time you run your application, profiling is available via websockets on port 6061,
// and once a client is connected it will send updates every 2 seconds. Defaults can be changed
// by passing one or more of the APIs options as params to Serve. See the examples for each option.
//
// To use the provided webserver, run the command "memstat" once your applications starts
// and has profiling enabled. To change HTTP port or connected to other sockets than default, see:
//	memstats --help
package memstats

import (
	"log"
	"net"
	"net/http"
	"runtime"
	"time"

	"golang.org/x/net/websocket"
)

type server struct {
	// ListenAddr is the address that the server listens on.
	ListenAddr string
	// Tick is the duration between two websocket updates.
	Tick time.Duration
	// MemRecordSize is the maximum number of records a profile will return.
	MemRecordSize int
}

func defaults(s *server) {
	s.ListenAddr = ":6061"
	s.Tick = 2 * time.Second
	s.MemRecordSize = 50
}

// Serve starts a memory monitoring server. By default it listens on :6061
func Serve(opts ...func(*server)) {
	var s server
	defaults(&s)
	for _, fn := range opts {
		fn(&s)
	}

	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		log.Fatalf("memstat: %s", err)
	}
	defer ln.Close()
	s.ListenAddr = ln.Addr().String()

	mux := http.NewServeMux()
	mux.Handle("/memstats-feed", websocket.Handler(s.ServeMemProfile))
	if err = http.Serve(ln, mux); err != nil {
		log.Fatalf("memstat: %s", err)
	}
}

// ServeMemProfile serves the connected socket with a snapshot of
// runtime.MemStats
func (s server) ServeMemProfile(ws *websocket.Conn) {
	defer ws.Close()
	var payload struct {
		runtime.MemStats
		Profile []memProfileRecord
		NumGo   int
	}
	for {
		if prof, ok := memProfile(s.MemRecordSize); ok {
			payload.Profile = prof
		}
		payload.NumGo = runtime.NumGoroutine()
		runtime.ReadMemStats(&payload.MemStats)
		err := websocket.JSON.Send(ws, payload)
		if err != nil {
			break
		}
		<-time.After(s.Tick)
	}
}

// memProfileRecord holds information about a memory profile entry
type memProfileRecord struct {
	runtime.MemProfileRecord
	// In use
	InUseObjs  int64
	InUseBytes int64
	// Stack trace
	Callstack []string
}

// memProfile returns a slice of memProfileRecord from the current memory profile.
func memProfile(size int) (data []memProfileRecord, ok bool) {
	record := make([]runtime.MemProfileRecord, size)
	n, ok := runtime.MemProfile(record, false)
	if !ok || n == 0 {
		return nil, false
	}
	prof := make([]memProfileRecord, len(record))
	for i, e := range record {
		prof[i] = memProfileRecord{
			MemProfileRecord: e,
			InUseBytes:       e.InUseBytes(),
			InUseObjs:        e.InUseObjects(),
			Callstack:        humanizeStack(e.Stack()),
		}
	}
	return prof[:n], true
}

// humanizeStack resolves a stracktrace to an array of function names
func humanizeStack(stk []uintptr) []string {
	fnpc := make([]string, len(stk))
	var n int
	for i, pc := range stk {
		fn := runtime.FuncForPC(pc)
		if fn == nil || pc == 0 {
			break
		}
		fnpc[i] = fn.Name()
		n++
	}
	return fnpc[:n]
}

// ListenAddr sets the address that the server will listen on for HTTP
// and WebSockets connections. ListenAddr is one of the options that can
// be provided to Serve.
func ListenAddr(addr string) func(*server) {
	return func(s *server) {
		s.ListenAddr = addr
	}
}

// Tick sets the frequency at which the websockets will send updates. Tick
// is one of the options that can be provided to Serve.
func Tick(d time.Duration) func(*server) {
	return func(s *server) {
		s.Tick = d
	}
}

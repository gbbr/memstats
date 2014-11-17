// Package memstats helps you monitor a running server's memory usage, visualize Garbage
// Collector information, run stack traces and memory profiles. The default values are
// configurable via the options provided by the API. To run the server, place this command
// at the top of your application:
//
// Example running with defaults (HTTP port :6061, refreshing every 2 seconds):
// 	go memstats.Serve()
// By default, the memory profile will be viewable on HTTP port :6061
package memstats

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/gbbr/memstats/internal/web"
	"golang.org/x/net/websocket"
)

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
	mux.Handle("/", s)
	mux.Handle("/memstats-feed", websocket.Handler(s.ServeMemProfile))
	if err = http.Serve(ln, mux); err != nil {
		log.Fatalf("memstat: %s", err)
	}
}

type server struct {
	ListenAddr    string
	Tick          time.Duration
	MemRecordSize int
}

// ServeHTTP serves the front-end HTML/JS viewer
func (s server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t, err := web.Template()
	if err != nil {
		fmt.Fprintf(w, "Error parsing template: %s", err)
		return
	}
	if err := t.ExecuteTemplate(w, "main", s); err != nil {
		fmt.Fprintf(w, "Error parsing template: %s", err)
	}
}

// ServeMemProfile serves the connected socket with a snapshot of
// runtime.MemStats
func (s server) ServeMemProfile(ws *websocket.Conn) {
	defer ws.Close()
	payload := struct {
		runtime.MemStats
		Profile []memProfileRecord
	}{}
	for {
		if prof, ok := memProfile(s.MemRecordSize); ok {
			payload.Profile = prof
		}
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
	AllocBytes, FreeBytes int64
	AllocObjs, FreeObjs   int64
	InUseBytes, InUseObjs int64
	Callstack             []string
}

// memProfile returns a slice of memProfileRecord from the current memory profile.
func memProfile(size int) (data []memProfileRecord, ok bool) {
	record := make([]runtime.MemProfileRecord, size)
	n, ok := runtime.MemProfile(record, false)
	if !ok {
		return nil, false
	}
	prof := make([]memProfileRecord, len(record))
	for i, e := range record {
		prof[i] = memProfileRecord{
			AllocBytes: e.AllocBytes,
			AllocObjs:  e.AllocObjects,
			FreeBytes:  e.FreeBytes,
			FreeObjs:   e.FreeObjects,
			InUseBytes: e.InUseBytes(),
			InUseObjs:  e.InUseObjects(),
			Callstack:  resolveFuncs(e.Stack()),
		}
	}
	return prof[:n], true
}

// resolveFuncs resolves a stracktrace to an array of function names
func resolveFuncs(stk []uintptr) []string {
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
// and WebSockets connections. The default port is :6061.
func ListenAddr(addr string) func(*server) {
	return func(s *server) {
		s.ListenAddr = addr
	}
}

// Tick sets the frequency at which the websockets will send updates.
// The default setting is 2 * time.Second.
func Tick(d time.Duration) func(*server) {
	return func(s *server) {
		s.Tick = d
	}
}

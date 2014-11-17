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
	mux.Handle("/", s)
	mux.Handle("/memstats-feed", websocket.Handler(s.ServeMemProfile))
	if err = http.Serve(ln, mux); err != nil {
		log.Fatalf("memstat: %s", err)
	}
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
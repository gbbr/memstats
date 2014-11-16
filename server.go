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

type profiler interface {
	profile() interface{}
}

type server struct {
	ListenAddr    string
	Tick          time.Duration
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
	mux.Handle("/memstats-feed", websocket.Handler(s.ServeMemStats))
	mux.Handle("/memprofile-feed", websocket.Handler(s.ServeMemProfile))
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

func (s server) serveProfiler(p profiler) websocket.Handler {
	return func(ws *websocket.Conn) {
		for {
			err := websocket.JSON.Send(ws, p.profile())
			if err != nil {
				break
			}
			<-time.After(s.Tick)
		}
	}
}

// ServeMemStats serves the connected socket with a snapshot of
// runtime.MemStats
func (s server) ServeMemStats(ws *websocket.Conn) {
	var buf runtime.MemStats
	for {
		runtime.ReadMemStats(&buf)
		err := websocket.JSON.Send(ws, buf)
		if err != nil {
			break
		}
		<-time.After(s.Tick)
	}
	ws.Close()
}

// ServeMemProfile serves the socket with memory profile blocks
func (s server) ServeMemProfile(ws *websocket.Conn) {
	var mp memProfile
	for {
		if data, ok := mp.profile(s.MemRecordSize); ok {
			err := websocket.JSON.Send(ws, data)
			if err != nil {
				break
			}
		}
		<-time.After(s.Tick)
	}
	ws.Close()
}

func ListenAddr(addr string) func(*server) {
	return func(s *server) {
		s.ListenAddr = addr
	}
}

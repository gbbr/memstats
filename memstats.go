package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"runtime"
	"time"

	"golang.org/x/net/websocket"
)

type server struct {
	ListenAddr string
	Tick       time.Duration
}

func Serve(opts ...func(*server)) {
	var m server

	defaults(&m)
	for _, fn := range opts {
		fn(&m)
	}

	ln, err := net.Listen("tcp", m.ListenAddr)
	if err != nil {
		log.Fatalf("memstat: %s", err)
	}
	m.ListenAddr = ln.Addr().String()

	mux := http.NewServeMux()
	mux.Handle("/memstats-feed", websocket.Handler(m.serveStats))
	mux.Handle("/", m)
	mux.Handle("/scripts/", http.FileServer(http.Dir("web")))

	if err = http.Serve(ln, mux); err != nil {
		log.Fatalf("memstat: %s", err)
	}
}

func (m server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("web/viewer.html")
	if err != nil {
		fmt.Fprintf(w, "Error parsing template: %s", err)
		return
	}
	if err := t.Execute(w, m); err != nil {
		fmt.Fprintf(w, "Error parsing template: %s", err)
	}
}

func (m server) serveStats(ws *websocket.Conn) {
	var stats runtime.MemStats
	for {
		runtime.ReadMemStats(&stats)
		websocket.JSON.Send(ws, stats)
		<-time.After(m.Tick)
	}
}

func defaults(s *server) {
	s.ListenAddr = ":6061"
	s.Tick = 2 * time.Second
}
func Addr(laddr string) func(*server) {
	return func(s *server) {
		s.ListenAddr = laddr
	}
}

func Duration(d time.Duration) func(*server) {
	return func(s *server) {
		s.Tick = d
	}
}

func main() {
	Serve()
}

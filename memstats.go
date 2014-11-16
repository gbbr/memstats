package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"time"

	"golang.org/x/net/websocket"
)

type server struct {
	ListenAddr string
	Tick       time.Duration
}

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
	s.ListenAddr = ln.Addr().String()

	mux := http.NewServeMux()
	mux.Handle("/", s)
	mux.Handle("/scripts/", http.FileServer(http.Dir("web")))
	mux.Handle("/memstats-feed", websocket.Handler(s.ServeSocket))
	if err = http.Serve(ln, mux); err != nil {
		log.Fatalf("memstat: %s", err)
	}
}

func (s server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("web/viewer.html")
	if err != nil {
		fmt.Fprintf(w, "Error parsing template: %s", err)
		return
	}
	if err := t.Execute(w, s); err != nil {
		fmt.Fprintf(w, "Error parsing template: %s", err)
	}
}

func (s server) ServeSocket(ws *websocket.Conn) {
	payload := struct {
		Stats runtime.MemStats
		CPU   string
	}{}
	var buf bytes.Buffer
	pprof.StartCPUProfile(&buf)
	for {
		runtime.ReadMemStats(&payload.Stats)
		payload.CPU = buf.String()
		buf.Reset()
		websocket.JSON.Send(ws, payload)
		<-time.After(s.Tick)
	}
	pprof.StopCPUProfile()
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
	go Serve()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}

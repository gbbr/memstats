package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"

	"golang.org/x/net/websocket"
)

type memStats struct {
	stats runtime.MemStats
	opts  config
}

type config struct {
	ListenAddr string
	Seconds    time.Duration
}

func Serve(addr string) {
	var m memStats

	m.opts.Seconds = 2 * time.Second

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("memstat: %s", err)
	}
	m.opts.ListenAddr = ln.Addr().String()

	mux := http.NewServeMux()
	mux.Handle("/memstats-feed", websocket.Handler(m.serveStats))
	mux.Handle("/", m)
	mux.Handle("/scripts/", http.FileServer(http.Dir("web")))

	if err = http.Serve(ln, mux); err != nil {
		log.Fatalf("memstat: %s", err)
	}
}

func (m memStats) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("web/viewer.html")
	if err != nil {
		fmt.Fprintf(w, "Error parsing template: %s", err)
		return
	}
	if err := t.Execute(w, m.opts); err != nil {
		fmt.Fprintf(w, "Error parsing template: %s", err)
	}
}

func (m memStats) serveStats(ws *websocket.Conn) {
	for {
		runtime.ReadMemStats(&m.stats)
		websocket.JSON.Send(ws, m.stats)

		<-time.After(m.opts.Seconds)
	}
}

func main() {
	Serve("127.0.0.1:8080")
}

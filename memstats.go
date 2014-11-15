package main

import (
	"fmt"
	"html/template"
	"log"
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

	m.opts.ListenAddr = "localhost:8080"
	m.opts.Seconds = 2 * time.Second

	mux := http.NewServeMux()
	mux.Handle("/memstats-feed", websocket.Handler(m.serveStats))
	mux.Handle("/", m)
	mux.Handle("/scripts/", http.FileServer(http.Dir("web")))

	err := http.ListenAndServe(m.opts.ListenAddr, mux)
	if err != nil {
		log.Fatalf("error starting MemStats server: %s", err)
	}
}

func (m memStats) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("web/viewer.html")
	if err != nil {
		fmt.Fprintf(w, "Error parsing template: %s", err)
		return
	}

	err = t.Execute(w, m.opts)
	if err != nil {
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
	Serve(":8080")
}

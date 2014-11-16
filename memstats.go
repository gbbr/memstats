package memstats

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/gbbr/memstats/internal/view"
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
	t, err := view.Render()
	if err != nil {
		fmt.Fprintf(w, "Error parsing template: %s", err)
		return
	}
	if err := t.ExecuteTemplate(w, "main", s); err != nil {
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

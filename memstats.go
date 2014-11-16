package memstats

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/gbbr/memstats/internal/web"
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
	mux.Handle("/memstats-feed", websocket.Handler(s.ServeMemStats))
	if err = http.Serve(ln, mux); err != nil {
		log.Fatalf("memstat: %s", err)
	}
}

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

func (s server) ServeMemStats(ws *websocket.Conn) {
	payload := struct {
		Stats runtime.MemStats
	}{}
	for {
		runtime.ReadMemStats(&payload.Stats)
		err := websocket.JSON.Send(ws, payload)
		if err != nil {
			break
		}
		<-time.After(s.Tick)
	}
	pprof.StopCPUProfile()
	ws.Close()
}

func defaults(s *server) {
	s.ListenAddr = ":6061"
	s.Tick = 2 * time.Second
}

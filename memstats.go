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
	ListenAddr      string
	Tick            time.Duration
	MemRecordSize   int
	BlockRecordSize int
}

func defaults(s *server) {
	s.ListenAddr = ":6061"
	s.Tick = 2 * time.Second
	s.MemRecordSize = 50
	s.BlockRecordSize = 50
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

type memProfile struct {
	AllocBytes, FreeBytes int64
	AllocObjs, FreeObjs   int64
	Callstack             []string
}

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

func payloadMem(p []runtime.MemProfileRecord) []memProfile {
	prof := make([]memProfile, len(p))
	for i, e := range p {
		prof[i] = memProfile{
			AllocBytes: e.AllocBytes,
			FreeBytes:  e.FreeBytes,
			AllocObjs:  e.AllocObjects,
			FreeObjs:   e.FreeObjects,
			Callstack:  resolveFuncs(e.Stack()),
		}
	}
	return prof
}

func (s server) ServeMemProfile(ws *websocket.Conn) {
	prof := make([]runtime.MemProfileRecord, s.MemRecordSize)
	for {
		n, ok := runtime.MemProfile(prof, false)
		if ok {
			err := websocket.JSON.Send(ws, payloadMem(prof[:n]))
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

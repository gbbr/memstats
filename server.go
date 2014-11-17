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

// ServeMemStats serves the connected socket with a snapshot of
// runtime.MemStats
func (s server) ServeMemStats(ws *websocket.Conn) {
	defer ws.Close()
	payload := struct {
		runtime.MemStats
		Profile []memProfile
	}{}
	for {
		if data, ok := runMemProfile(s.MemRecordSize); ok {
			payload.Profile = data
		}
		runtime.ReadMemStats(&payload.MemStats)
		err := websocket.JSON.Send(ws, payload)
		if err != nil {
			break
		}
		<-time.After(s.Tick)
	}
}

// memProfile holds information about a memory profile entry
type memProfile struct {
	AllocBytes, FreeBytes int64
	AllocObjs, FreeObjs   int64
	InUseBytes, InUseObjs int64
	Callstack             []string
}

func runMemProfile(size int) (data []memProfile, ok bool) {
	record := make([]runtime.MemProfileRecord, size)
	n, ok := runtime.MemProfile(record, false)
	if !ok {
		return nil, false
	}
	prof := make([]memProfile, len(record))
	for i, e := range record {
		prof[i] = memProfile{
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

func ListenAddr(addr string) func(*server) {
	return func(s *server) {
		s.ListenAddr = addr
	}
}

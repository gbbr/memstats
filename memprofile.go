package memstats

import "runtime"

// memProfile holds information about a memory profile entry
type memProfile struct {
	AllocBytes, FreeBytes int64
	AllocObjs, FreeObjs   int64
	InUseBytes, InUseObjs int64
	Callstack             []string
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

func (m memProfile) payload(size int) (data []memProfile, ok bool) {
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

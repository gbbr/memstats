package memstats

import "runtime"

// memProfileRecord holds information about a memory profile entry
type memProfileRecord struct {
	// Objects
	AllocObjs int64
	FreeObjs  int64
	InUseObjs int64
	// Byte values
	AllocBytes int64
	FreeBytes  int64
	InUseBytes int64
	// Stack trace
	Callstack []string
}

// memProfile returns a slice of memProfileRecord from the current memory profile.
func memProfile(size int) (data []memProfileRecord, ok bool) {
	record := make([]runtime.MemProfileRecord, size)
	n, ok := runtime.MemProfile(record, false)
	if !ok || n == 0 {
		return nil, false
	}
	prof := make([]memProfileRecord, len(record))
	for i, e := range record {
		prof[i] = memProfileRecord{
			// Bytes
			AllocBytes: e.AllocBytes,
			FreeBytes:  e.FreeBytes,
			InUseBytes: e.InUseBytes(),
			// Objects
			AllocObjs: e.AllocObjects,
			FreeObjs:  e.FreeObjects,
			InUseObjs: e.InUseObjects(),
			// Stack
			Callstack: resolveFuncs(e.Stack()),
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

package syncmap

import (
	"debug/dwarf"
	"sync"
	"sync/atomic"
)

type Map struct {
	counter *int64
	mu      sync.Mutex
	read    atomic.Value //readOnly
	dirty   map[interface{}]*dwarf.Entry
	misses  int
}

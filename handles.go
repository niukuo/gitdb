package gitdb

import (
	"fmt"
	"sync"
	"unsafe"
)

type HandleList struct {
	sync.RWMutex
	// stores the Go pointers
	handles map[unsafe.Pointer]interface{}
}

func NewHandleList() *HandleList {
	return &HandleList{
		handles: make(map[unsafe.Pointer]interface{}),
	}
}

func (l *HandleList) Track(handle unsafe.Pointer, value interface{}) {
	l.Lock()
	defer l.Unlock()
	if _, ok := l.handles[handle]; ok {
		panic(fmt.Sprintf("handle already tracked: %p", handle))
	}
	l.handles[handle] = value
}

func (l *HandleList) Untrack(handle unsafe.Pointer) interface{} {
	l.Lock()
	defer l.Unlock()
	v, ok := l.handles[handle]
	if !ok {
		panic(fmt.Sprintf("untracking invalid handle: %p", handle))
	}
	delete(l.handles, handle)
	return v
}

func (l *HandleList) Get(handle unsafe.Pointer) interface{} {
	l.RLock()
	defer l.RUnlock()

	ptr, ok := l.handles[handle]
	if !ok {
		panic(fmt.Sprintf("invalid pointer handle: %p", handle))
	}

	return ptr
}

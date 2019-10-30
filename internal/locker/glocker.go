package locker

import "sync"

type mutex struct {
	sync.Mutex
	ref int32
}

//GLocker 本地锁
type GLocker struct {
	mux     sync.Mutex
	lockers map[string]*mutex
}

//NewGLocker 新建本地锁
func NewGLocker() NamedLocker {
	return &GLocker{
		lockers: make(map[string]*mutex),
	}
}

//Lock Locker.Lock的实现
func (g *GLocker) Lock(name string) {
	g.mux.Lock()
	var mux *mutex
	var ok bool
	if mux, ok = g.lockers[name]; !ok {
		mux = &mutex{}
		g.lockers[name] = mux
	}
	mux.ref++
	g.mux.Unlock()
	mux.Lock()
}

//UnLock Locker.UnLock
func (g *GLocker) UnLock(name string) {
	g.mux.Lock()
	defer g.mux.Unlock()

	if mux, ok := g.lockers[name]; ok {
		mux.Unlock()
		mux.ref--
		if mux.ref == 0 {
			delete(g.lockers, name)
		}
	}
}

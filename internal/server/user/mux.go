package user

import (
	"sync"

	"github.com/gorilla/mux"
)

type swappableMux struct {
	m    sync.RWMutex
	root mux.Router
}

func (sm *swappableMux) get() mux.Router {
	sm.m.RLock()
	defer sm.m.RUnlock()

	return sm.root
}

func (sm *swappableMux) set(mux mux.Router) {
	sm.m.Lock()
	sm.root = mux
	sm.m.Unlock()
}
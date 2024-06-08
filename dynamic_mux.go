package grimoire

import (
	"net/http"
	"sync"
)

type dynamicMux struct {
	mux *http.ServeMux
	sync.RWMutex
}

func (d *dynamicMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.RLock()
	currentMux := d.mux
	d.RUnlock()
	currentMux.ServeHTTP(w, r)
}

func (d *dynamicMux) Update(mux *http.ServeMux) {
	d.Lock()
	d.mux = mux
	d.Unlock()
}

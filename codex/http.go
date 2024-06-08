package codex

import (
	"context"
	"encoding/json"
	"fmt"
	g "github.com/ananrafs/grimoire"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type HttpCodex struct {
	routes         []g.Route
	routeLock      sync.Mutex
	onChangeSignal chan struct{}

	server *http.Server
	quit   chan bool
}

func NewHttpCodex(port int) (out *HttpCodex, onQuit func()) {
	_default := &HttpCodex{
		onChangeSignal: make(chan struct{}),
		quit:           make(chan bool),
		routes:         make([]g.Route, 0),
		routeLock:      sync.Mutex{},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/addRoute", func(w http.ResponseWriter, r *http.Request) {
		var route g.Route
		err := json.NewDecoder(r.Body).Decode(&route)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_default.routeLock.Lock()
		defer _default.routeLock.Unlock()
		_default.routes = append(_default.routes, route)

		_default.onChangeSignal <- struct{}{}

		w.WriteHeader(http.StatusOK)
	})

	_default.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	return _default, func() {
		_default.quit <- true

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := _default.server.Shutdown(ctx); err != nil {
			log.Fatalf("Server Shutdown: %v", err)
		}

		os.Exit(0)
	}
}

func (hc *HttpCodex) Init() error {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := hc.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	go func() {
		for _ = range quit {
			hc.quit <- true
		}
	}()

	return nil
}

func (hc *HttpCodex) GetAllRoute() []g.Route {
	hc.routeLock.Lock()
	defer hc.routeLock.Unlock()
	return hc.routes
}

func (hc *HttpCodex) GetChannel() chan struct{} {
	return hc.onChangeSignal
}

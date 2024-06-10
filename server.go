package grimoire

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Server interface {
	Serve(port string) (onQuit func())
}

type serverOpts func(server *servant)

type servant struct {
	codex  Codex
	caster Caster[http.Request]
	dMux   *dynamicMux
	mux    sync.RWMutex

	onQuit func()

	signalHandler func(sign <-chan struct{}, act func())
	logger        logger
}

func NewServer(codex Codex, caster Caster[http.Request], opts ...serverOpts) Server {
	_default := &servant{
		codex:         codex,
		caster:        caster,
		dMux:          &dynamicMux{mux: http.NewServeMux()},
		signalHandler: defaultSignalHandler,
		logger:        newdefaultLogger(),
	}
	for _, opt := range opts {
		opt(_default)
	}
	return _default
}

func (s *servant) Serve(port string) (onQuit func()) {
	if err := s.codex.Init(); err != nil {
		s.logger.Error(fmt.Errorf("error initializing codex: %v", err))
	}

	go s.signalHandler(s.codex.GetChannel(), s.refreshHandlers)

	s.refreshHandlers()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: s.dMux,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error(fmt.Errorf("listenAndServe(): %v", err))
		}
	}()

	return func() {
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			s.logger.Error(fmt.Errorf("server Shutdown: %v", err))
		}
		s.logger.Debug("Server gracefully stopped")
	}

}

func (s *servant) refreshHandlers() {
	s.mux.Lock()
	defer s.mux.Unlock()

	newMux := http.NewServeMux()

	for _, route := range s.codex.GetAllRoute() {
		newMux.HandleFunc(route.URL,
			s.generateHandlerFunc(route.Meta),
		)

	}

	s.dMux.Update(newMux)
}

func (s *servant) generateHandlerFunc(meta []Meta) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := s.caster.Cast(meta, *r)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// Write the response to the client
		json.NewEncoder(w).Encode(response)
	}
}

func defaultSignalHandler(signal <-chan struct{}, act func()) {
	for _ = range signal {
		act()
	}
}

func WithThrottle(duration time.Duration) serverOpts {
	return func(server *servant) {
		server.signalHandler = func(signal <-chan struct{}, act func()) {
			throttle(duration, signal, act)
		}
	}
}

func WithCustomLogger(logger logger) serverOpts {
	return func(server *servant) {
		server.logger = logger
	}
}

package router

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/raflynagachi/go-rest-api-starter/config"
	hn "github.com/raflynagachi/go-rest-api-starter/internal/handler/definition"
	"github.com/raflynagachi/go-rest-api-starter/pkg/logger"
)

type Router struct {
	Cfg       *config.Config
	Router    *httprouter.Router
	appLogger *logger.Logger
	server    *http.Server
	mu        sync.Mutex // mutex to ensure thread-safe access
}

// New creates a new Router instance
func New(cfg *config.Config, log *logger.Logger, hn hn.APIHandler) *Router {
	router := newRouter(hn)
	return &Router{
		Cfg:       cfg,
		appLogger: log,
		Router:    router,
	}
}

// Start initializes and starts the HTTP server
func (r *Router) Start() error {
	addr := fmt.Sprintf(":%d", r.Cfg.App.Port)
	r.server = &http.Server{
		Addr:    addr,
		Handler: r.Router,
	}

	r.appLogger.Info(fmt.Sprintf("Running on %s", addr))
	return r.server.ListenAndServe()
}

// Shutdown gracefully stops the HTTP server
func (r *Router) Shutdown(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.server == nil {
		return nil // server was not started
	}

	r.appLogger.Info("Shutting down server")
	return r.server.Shutdown(ctx)
}

// ServeHTTP handles the HTTP server's lifecycle
func (r *Router) ServeHTTP() error {
	if err := r.Start(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

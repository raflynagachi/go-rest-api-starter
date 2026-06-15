package router

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"github.com/raflynagachi/go-rest-api-starter/config"
	hn "github.com/raflynagachi/go-rest-api-starter/internal/payment/handler/definition"
	"github.com/raflynagachi/go-rest-api-starter/internal/handler/middleware"
	"github.com/raflynagachi/go-rest-api-starter/pkg/logger"
)

type Router struct {
	Cfg       *config.Config
	Router    *httprouter.Router
	appLogger *logger.Logger
	server    *http.Server
	mu        sync.Mutex
}

func New(cfg *config.Config, log *logger.Logger, hn hn.APIHandler, db *sqlx.DB) *Router {
	return &Router{
		Cfg:       cfg,
		appLogger: log,
		Router:    newRouter(cfg, log, hn, db),
	}
}

func (r *Router) Start() error {
	addr := fmt.Sprintf(":%d", r.Cfg.App.Port)
	r.server = &http.Server{
		Addr:    addr,
		Handler: middleware.CORS(middleware.Logging(r.appLogger, middleware.RequestID(r.Router))),
	}
	r.appLogger.Info(fmt.Sprintf("Running on %s", addr))
	return r.server.ListenAndServe()
}

func (r *Router) Shutdown(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.server == nil {
		return nil
	}
	r.appLogger.Info("Shutting down server")
	return r.server.Shutdown(ctx)
}

func (r *Router) ServeHTTP() error {
	if err := r.Start(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

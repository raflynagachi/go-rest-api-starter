package router

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"github.com/raflynagachi/go-rest-api-starter/config"
	hn "github.com/raflynagachi/go-rest-api-starter/internal/payment/handler/definition"
	"github.com/raflynagachi/go-rest-api-starter/internal/handler/middleware"
	"github.com/raflynagachi/go-rest-api-starter/pkg/logger"
)

func newRouter(cfg *config.Config, log *logger.Logger, hn hn.APIHandler, db *sqlx.DB) *httprouter.Router {
	router := httprouter.New()

	jwtAuth := func(h httprouter.Handle) httprouter.Handle {
		return middleware.JWTAuth(cfg.JwtKey, log, h)
	}

	router.GET("/ping", Ping)
	router.GET("/health", Health(db))
	router.GET("/transactions/:id", hn.GetTransactionByID)
	router.POST("/transactions", jwtAuth(hn.CreateTransaction))

	return router
}

func Ping(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func Health(db *sqlx.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		if err := db.PingContext(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{"status": "unavailable"})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}

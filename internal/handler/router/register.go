package router

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"github.com/raflynagachi/go-rest-api-starter/config"
	hn "github.com/raflynagachi/go-rest-api-starter/internal/handler/definition"
	"github.com/raflynagachi/go-rest-api-starter/internal/handler/middleware"
	"github.com/raflynagachi/go-rest-api-starter/pkg/logger"
)

func newRouter(cfg *config.Config, log *logger.Logger, hn hn.APIHandler, db *sqlx.DB) *httprouter.Router {
	router := httprouter.New()

	jwtAuth := func(h httprouter.Handle) httprouter.Handle {
		return middleware.JWTAuth(cfg.JwtKey, log, h)
	}

	// public routes
	router.GET("/ping", Ping)
	router.GET("/health", Health(db))
	router.POST("/auth/register", hn.Register)
	router.POST("/auth/login", hn.Login)
	router.GET("/users", hn.GetUser)
	router.GET("/users/:id", hn.GetUserByID)

	// protected routes (require JWT Bearer token)
	router.POST("/users", jwtAuth(hn.CreateUser))
	router.PUT("/users/:id", jwtAuth(hn.UpdateUser))
	router.DELETE("/users/:id", jwtAuth(hn.DeleteUser))

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

package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	hn "github.com/raflynagachi/go-rest-api-starter/internal/handler/definition"
)

func newRouter(hn hn.APIHandler) *httprouter.Router {
	router := httprouter.New()
	// middlewares

	// API
	router.GET("/users", hn.GetUser)
	router.GET("/users/:id", hn.GetUserByID)
	router.POST("/users", hn.CreateUser)
	router.PUT("/users/:id", hn.UpdateUser)

	router.GET("/ping", Ping)

	return router
}

func Ping(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

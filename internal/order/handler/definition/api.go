package definition

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type APIHandler interface {
	GetOrders(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetOrderByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	CreateOrder(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

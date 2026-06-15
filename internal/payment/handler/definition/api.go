package definition

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type APIHandler interface {
	GetTransactionByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	CreateTransaction(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

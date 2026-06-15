package definition

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type APIHandler interface {
	GetItems(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetItemByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	CreateItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	UpdateStock(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

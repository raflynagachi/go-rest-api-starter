package definition

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type APIHandler interface {
	GetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetUserByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	CreateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	UpdateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

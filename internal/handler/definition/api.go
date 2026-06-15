package definition

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type APIHandler interface {
	Register(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	GetUserByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	CreateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	UpdateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	DeleteUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

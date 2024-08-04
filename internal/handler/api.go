package handler

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	req "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/request"
	"github.com/raflynagachi/go-rest-api-starter/pkg/http/encoder"
	"github.com/raflynagachi/go-rest-api-starter/pkg/http/response"
)

func (h *APIHandlerImpl) GetUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	filter := req.UserFilter{}
	err := populateStructFromQueryParams(r, &filter)
	if err != nil {
		response.WriteFromError(w, response.WrapErrBadRequest(err))
		return
	}

	ctx := r.Context()
	resp, err := h.usecase.GetUser(ctx, filter)
	if err != nil {
		response.WriteFromError(w, err)
		return
	}

	response.WriteOKResponse(w, resp)
}

func (h *APIHandlerImpl) GetUserByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		response.WriteFromError(w, response.WrapErrBadRequest(err))
		return
	}

	ctx := r.Context()
	resp, err := h.usecase.GetUserByID(ctx, int64(id))
	if err != nil {
		response.WriteFromError(w, err)
		return
	}

	response.WriteOKResponse(w, resp)
}

func (a *APIHandlerImpl) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := &req.CreateUpdateUserReq{}
	err := encoder.DecodeJson(r, req)
	if err != nil {
		response.WriteFromError(w, response.WrapErrBadRequest(err))
		return
	}

	err = a.usecase.CreateUser(r.Context(), req)
	if err != nil {
		response.WriteFromError(w, err)
		return
	}

	response.WriteOKResponse(w, "create User success")
}

func (a *APIHandlerImpl) UpdateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		response.WriteFromError(w, response.WrapErrBadRequest(err))
		return
	}

	req := &req.CreateUpdateUserReq{}
	err = encoder.DecodeJson(r, req)
	if err != nil {
		response.WriteFromError(w, response.WrapErrBadRequest(err))
		return
	}

	err = a.usecase.UpdateUser(r.Context(), int64(id), req)
	if err != nil {
		response.WriteFromError(w, err)
		return
	}

	response.WriteOKResponse(w, "update User success")
}

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
		response.WriteFromError(w, r, response.WrapErrBadRequest(err), h.appLogger)
		return
	}

	ctx := r.Context()
	resp, err := h.usecase.GetUser(ctx, filter)
	if err != nil {
		response.WriteFromError(w, r, err, h.appLogger)
		return
	}

	response.WriteOKResponse(w, r, resp, h.appLogger)
}

func (h *APIHandlerImpl) GetUserByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		response.WriteFromError(w, r, response.WrapErrBadRequest(err), h.appLogger)
		return
	}

	ctx := r.Context()
	resp, err := h.usecase.GetUserByID(ctx, int64(id))
	if err != nil {
		response.WriteFromError(w, r, err, h.appLogger)
		return
	}

	response.WriteOKResponse(w, r, resp, h.appLogger)
}

func (h *APIHandlerImpl) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := &req.CreateUpdateUserReq{}
	err := encoder.DecodeJson(r, req)
	if err != nil {
		response.WriteFromError(w, r, response.WrapErrBadRequest(err), h.appLogger)
		return
	}

	err = h.usecase.CreateUser(r.Context(), req)
	if err != nil {
		response.WriteFromError(w, r, err, h.appLogger)
		return
	}

	response.WriteOKResponse(w, r, "create User success", h.appLogger)
}

func (h *APIHandlerImpl) UpdateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		response.WriteFromError(w, r, response.WrapErrBadRequest(err), h.appLogger)
		return
	}

	req := &req.CreateUpdateUserReq{}
	err = encoder.DecodeJson(r, req)
	if err != nil {
		response.WriteFromError(w, r, response.WrapErrBadRequest(err), h.appLogger)
		return
	}

	err = h.usecase.UpdateUser(r.Context(), int64(id), req)
	if err != nil {
		response.WriteFromError(w, r, err, h.appLogger)
		return
	}

	response.WriteOKResponse(w, r, "update User success", h.appLogger)
}

package handler

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	req "github.com/raflynagachi/go-rest-api-starter/internal/inventory/dto/web/request"
	"github.com/raflynagachi/go-rest-api-starter/pkg/http/encoder"
	"github.com/raflynagachi/go-rest-api-starter/pkg/http/response"
)

func (h *APIHandlerImpl) GetItems(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	filter := req.ItemFilter{}
	if err := populateStructFromQueryParams(r, &filter); err != nil {
		response.WriteFromError(w, r, response.WrapErrBadRequest(err), h.appLogger)
		return
	}

	res, err := h.usecase.GetItems(r.Context(), filter)
	if err != nil {
		response.WriteFromError(w, r, err, h.appLogger)
		return
	}
	response.WriteOKResponse(w, r, res, h.appLogger)
}

func (h *APIHandlerImpl) GetItemByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		response.WriteFromError(w, r, response.WrapErrBadRequest(err), h.appLogger)
		return
	}

	res, err := h.usecase.GetItemByID(r.Context(), id)
	if err != nil {
		response.WriteFromError(w, r, err, h.appLogger)
		return
	}
	response.WriteOKResponse(w, r, res, h.appLogger)
}

func (h *APIHandlerImpl) CreateItem(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body := &req.CreateItemReq{}
	if err := encoder.DecodeJson(r, body); err != nil {
		response.WriteFromError(w, r, response.WrapErrBadRequest(err), h.appLogger)
		return
	}

	if err := h.usecase.CreateItem(r.Context(), body); err != nil {
		response.WriteFromError(w, r, err, h.appLogger)
		return
	}
	response.WriteOKResponse(w, r, "item created", h.appLogger)
}

func (h *APIHandlerImpl) UpdateStock(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		response.WriteFromError(w, r, response.WrapErrBadRequest(err), h.appLogger)
		return
	}

	body := &req.UpdateStockReq{}
	if err := encoder.DecodeJson(r, body); err != nil {
		response.WriteFromError(w, r, response.WrapErrBadRequest(err), h.appLogger)
		return
	}

	if err := h.usecase.UpdateStock(r.Context(), id, body); err != nil {
		response.WriteFromError(w, r, err, h.appLogger)
		return
	}
	response.WriteOKResponse(w, r, "stock updated", h.appLogger)
}

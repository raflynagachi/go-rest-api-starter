package handler

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	req "github.com/raflynagachi/go-rest-api-starter/internal/order/dto/web/request"
	"github.com/raflynagachi/go-rest-api-starter/pkg/http/encoder"
	"github.com/raflynagachi/go-rest-api-starter/pkg/http/response"
)

func (h *APIHandlerImpl) GetOrders(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	filter := req.OrderFilter{}
	if err := populateStructFromQueryParams(r, &filter); err != nil {
		response.WriteFromError(w, r, response.WrapErrBadRequest(err), h.appLogger)
		return
	}
	res, err := h.usecase.GetOrders(r.Context(), filter)
	if err != nil {
		response.WriteFromError(w, r, err, h.appLogger)
		return
	}
	response.WriteOKResponse(w, r, res, h.appLogger)
}

func (h *APIHandlerImpl) GetOrderByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		response.WriteFromError(w, r, response.WrapErrBadRequest(err), h.appLogger)
		return
	}
	res, err := h.usecase.GetOrderByID(r.Context(), id)
	if err != nil {
		response.WriteFromError(w, r, err, h.appLogger)
		return
	}
	response.WriteOKResponse(w, r, res, h.appLogger)
}

func (h *APIHandlerImpl) CreateOrder(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body := &req.CreateOrderReq{}
	if err := encoder.DecodeJson(r, body); err != nil {
		response.WriteFromError(w, r, response.WrapErrBadRequest(err), h.appLogger)
		return
	}
	res, err := h.usecase.CreateOrder(r.Context(), body)
	if err != nil {
		response.WriteFromError(w, r, err, h.appLogger)
		return
	}
	response.WriteOKResponse(w, r, res, h.appLogger)
}

package handler

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	req "github.com/raflynagachi/go-rest-api-starter/internal/payment/dto/web/request"
	"github.com/raflynagachi/go-rest-api-starter/pkg/http/encoder"
	"github.com/raflynagachi/go-rest-api-starter/pkg/http/response"
)

func (h *APIHandlerImpl) GetTransactionByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		response.WriteFromError(w, r, response.WrapErrBadRequest(err), h.appLogger)
		return
	}
	res, err := h.usecase.GetTransactionByID(r.Context(), id)
	if err != nil {
		response.WriteFromError(w, r, err, h.appLogger)
		return
	}
	response.WriteOKResponse(w, r, res, h.appLogger)
}

func (h *APIHandlerImpl) CreateTransaction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body := &req.CreateTransactionReq{}
	if err := encoder.DecodeJson(r, body); err != nil {
		response.WriteFromError(w, r, response.WrapErrBadRequest(err), h.appLogger)
		return
	}
	res, err := h.usecase.CreateTransaction(r.Context(), body)
	if err != nil {
		response.WriteFromError(w, r, err, h.appLogger)
		return
	}
	response.WriteOKResponse(w, r, res, h.appLogger)
}

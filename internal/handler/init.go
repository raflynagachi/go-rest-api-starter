package handler

import (
	hn "github.com/raflynagachi/go-rest-api-starter/internal/handler/definition"
	uc "github.com/raflynagachi/go-rest-api-starter/internal/usecase/definition"
	"github.com/raflynagachi/go-rest-api-starter/pkg/http/request"
)

type APIHandlerImpl struct {
	usecase uc.APIUsecase
}

func New(usecase uc.APIUsecase) hn.APIHandler {
	return &APIHandlerImpl{
		usecase: usecase,
	}
}

var (
	populateStructFromQueryParams = request.PopulateStructFromQueryParams
)

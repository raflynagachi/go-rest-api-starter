package handler

import (
	hn "github.com/raflynagachi/go-rest-api-starter/internal/handler/definition"
	uc "github.com/raflynagachi/go-rest-api-starter/internal/usecase/definition"
	"github.com/raflynagachi/go-rest-api-starter/pkg/http/request"
	"github.com/raflynagachi/go-rest-api-starter/pkg/logger"
)

type APIHandlerImpl struct {
	usecase   uc.APIUsecase
	appLogger *logger.Logger
}

func New(usecase uc.APIUsecase, log *logger.Logger) hn.APIHandler {
	return &APIHandlerImpl{
		usecase:   usecase,
		appLogger: log,
	}
}

var (
	populateStructFromQueryParams = request.PopulateStructFromQueryParams
)

package response

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/raflynagachi/go-rest-api-starter/pkg/http/encoder"
	appValidator "github.com/raflynagachi/go-rest-api-starter/pkg/validator"
	"github.com/rs/zerolog/log"
)

const ()

type ErrResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e ErrResponse) Error() string {
	return e.Err.Error()
}

func (e ErrResponse) Cause() error {
	return e
}

// WriteFromError writes a formatted error response
func WriteFromError(w http.ResponseWriter, e error) {
	errResp, lastError := FindErrResponse(e)
	errResp.Message = lastError.Error()

	if errResp.Code == http.StatusBadRequest {
		if valErrs, ok := errResp.Err.(validator.ValidationErrors); ok {
			var msgBuilder strings.Builder
			for i, fieldErr := range valErrs {
				msgBuilder.WriteString(fieldErr.Translate(appValidator.GetTranslator()))
				if i < len(valErrs)-1 {
					msgBuilder.WriteString(",")
				}
			}
			errResp.Message = msgBuilder.String()
		}
	}

	log.Error().Msg(fmt.Sprintf("%+v with status code %d", errResp, errResp.Code))
	if errResp.Code == http.StatusInternalServerError {
		errResp.Err = errors.New("internal server error")
		errResp.Message = errResp.Err.Error()
	}

	w.WriteHeader(errResp.Code)
	err := encoder.EncodeJson(w, errResp)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func WrapErrBadRequest(err error) ErrResponse {
	return ErrResponse{
		Code: http.StatusBadRequest,
		Err:  err,
	}
}

func WrapErrNotFound(err error) ErrResponse {
	return ErrResponse{
		Code: http.StatusNotFound,
		Err:  err,
	}
}

func WrapErrInternalServer(err error) ErrResponse {
	return ErrResponse{
		Code: http.StatusInternalServerError,
		Err:  err,
	}
}

// FindErrResponse recursively finds the root ErrResponse in the error chain
func FindErrResponse(err error) (errResp ErrResponse, lastError error) {
	for err != nil {
		if errRespTmp, ok := err.(ErrResponse); ok {
			errResp = errRespTmp
			break
		}
		err = errors.Unwrap(err)
	}

	lastError = errors.Cause(errResp.Err)

	if (errResp == ErrResponse{}) {
		return WrapErrInternalServer(err), lastError
	}

	return errResp, lastError
}

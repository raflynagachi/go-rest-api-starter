package response

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	appValidator "github.com/raflynagachi/go-rest-api-starter/pkg/validator"
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

// WriteFromError writes a formatted error response
func WriteFromError(w http.ResponseWriter, r *http.Request, e error, log *slog.Logger) {
	errResp, lastError := findErrResponse(e)
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

	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	body := string(bodyBytes)

	log.Error(
		"error response",
		slog.Int("code", errResp.Code),
		slog.String("path", r.URL.Path),
		slog.String("method", r.Method),
		slog.String("request", body),
		slog.String("error", errResp.Error()),
	)

	if errResp.Code == http.StatusInternalServerError {
		errResp.Err = errors.New("internal server error")
		errResp.Message = errResp.Err.Error()
	}

	w.WriteHeader(errResp.Code)
	err := encodeJson(w, errResp)
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

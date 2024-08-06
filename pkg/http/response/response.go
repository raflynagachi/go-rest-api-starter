package response

import (
	"log/slog"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Meta `json:"-"`
}

func WriteResponse(w http.ResponseWriter, response Response, log *slog.Logger) {
	log.Info(
		"send response",
		slog.Int("code", response.Code),
		slog.String("path", response.Meta.Path),
		slog.String("method", response.Meta.Method),
	)

	w.WriteHeader(response.Code)
	err := encodeJson(w, response)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func WriteOKResponse(w http.ResponseWriter, r *http.Request, data interface{}, logger *slog.Logger) {
	WriteResponse(w, Response{
		Code: http.StatusOK,
		Data: data,
		Meta: Meta{
			Path:   r.URL.Path,
			Method: r.Method,
		},
	}, logger)
}

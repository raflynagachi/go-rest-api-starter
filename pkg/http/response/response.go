package response

import (
	"fmt"
	"net/http"

	"github.com/raflynagachi/go-rest-api-starter/pkg/http/encoder"
	"github.com/rs/zerolog/log"
)

var (
	encodeJson = encoder.EncodeJson
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta,omitempty"`
}

func WriteResponse(w http.ResponseWriter, response Response) {
	log.Info().Msg(fmt.Sprintf("%+v", response))

	w.WriteHeader(response.Code)
	err := encodeJson(w, response)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func WriteOKResponse(w http.ResponseWriter, data interface{}) {
	WriteResponse(w, Response{
		Code: http.StatusOK,
		Data: data,
	})
}

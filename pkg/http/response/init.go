package response

import (
	"io"
	"log/slog"

	"github.com/raflynagachi/go-rest-api-starter/pkg/http/encoder"
)

var (
	mockLogger = slog.New(slog.NewJSONHandler(io.Discard, nil))

	encodeJson      = encoder.EncodeJson
	findErrResponse = FindErrResponse
)

type Meta struct {
	Path   string
	Method string
}

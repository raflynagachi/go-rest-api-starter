package encoder

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func EncodeJson(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(data)
}

func DecodeJson(r *http.Request, data interface{}) error {
	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	body := io.NopCloser(bytes.NewBuffer(bodyBytes))
	return json.NewDecoder(body).Decode(data)
}

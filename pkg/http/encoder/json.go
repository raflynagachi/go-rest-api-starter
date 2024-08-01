package encoder

import (
	"encoding/json"
	"net/http"
)

func EncodeJson(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(data)
}

func DecodeJson(r *http.Request, data interface{}) error {
	r.Header.Set("Content-Type", "application/json")
	return json.NewDecoder(r.Body).Decode(data)
}
